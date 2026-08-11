package builder

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The EPEL keys have to survive the unmarshal/re-marshal the CLI does before it
// hands the config to Packer as a var file. An explicit empty value is what
// disables EPEL, so it has to reach Packer; a config that never mentions EPEL
// must not gain the keys, or every existing RHEL build would lose EPEL.
func TestBaremetalConfigEpelRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want map[string]interface{}
	}{
		{
			name: "explicit empty values are preserved to disable EPEL",
			in:   `{"disk_size":"20480","redhat_epel_rpm":"","epel_rpm_gpg_key":""}`,
			want: map[string]interface{}{
				"redhat_epel_rpm":  "",
				"epel_rpm_gpg_key": "",
			},
		},
		{
			name: "custom mirrors are preserved",
			in:   `{"redhat_epel_rpm":"https://mirror.internal/epel-release-latest-9.noarch.rpm","epel_rpm_gpg_key":"https://mirror.internal/RPM-GPG-KEY-EPEL-9"}`,
			want: map[string]interface{}{
				"redhat_epel_rpm":  "https://mirror.internal/epel-release-latest-9.noarch.rpm",
				"epel_rpm_gpg_key": "https://mirror.internal/RPM-GPG-KEY-EPEL-9",
			},
		},
		{
			name: "only one key set leaves the other to the upstream default",
			in:   `{"redhat_epel_rpm":""}`,
			want: map[string]interface{}{
				"redhat_epel_rpm": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config BaremetalConfig
			assert.NoError(t, json.Unmarshal([]byte(tt.in), &config))

			out, err := json.Marshal(config)
			assert.NoError(t, err)

			var got map[string]interface{}
			assert.NoError(t, json.Unmarshal(out, &got))

			for key, want := range tt.want {
				value, present := got[key]
				assert.Truef(t, present, "%s must reach Packer", key)
				assert.Equal(t, want, value)
			}
		})
	}
}

// Guards the regression that a non-pointer field would cause: configs that say
// nothing about EPEL must not emit the keys at all, otherwise the empty value
// would disable EPEL for builds that never asked.
func TestConfigsOmitEpelWhenUnset(t *testing.T) {
	tests := []struct {
		name   string
		config interface{}
	}{
		{name: "baremetal", config: &BaremetalConfig{}},
		{name: "vsphere", config: &VsphereConfig{}},
		{name: "cloudstack", config: &CloudstackConfig{}},
		{name: "nutanix", config: &NutanixConfig{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, json.Unmarshal([]byte(`{}`), tt.config))

			out, err := json.Marshal(tt.config)
			assert.NoError(t, err)

			var got map[string]interface{}
			assert.NoError(t, json.Unmarshal(out, &got))

			assert.NotContains(t, got, "redhat_epel_rpm")
			assert.NotContains(t, got, "epel_rpm_gpg_key")
		})
	}
}

// Every hypervisor config that embeds RhelConfig has the same defect, so the
// fix has to cover all of them.
func TestEpelConfigurableOnAllRhelHypervisors(t *testing.T) {
	in := []byte(`{"redhat_epel_rpm":"","epel_rpm_gpg_key":""}`)

	assertDisabled := func(t *testing.T, config interface{}) {
		t.Helper()
		assert.NoError(t, json.Unmarshal(in, config))

		out, err := json.Marshal(config)
		assert.NoError(t, err)

		var got map[string]interface{}
		assert.NoError(t, json.Unmarshal(out, &got))
		assert.Equal(t, "", got["redhat_epel_rpm"])
		assert.Equal(t, "", got["epel_rpm_gpg_key"])
	}

	t.Run("baremetal", func(t *testing.T) { assertDisabled(t, &BaremetalConfig{}) })
	t.Run("vsphere", func(t *testing.T) { assertDisabled(t, &VsphereConfig{}) })
	t.Run("cloudstack", func(t *testing.T) { assertDisabled(t, &CloudstackConfig{}) })
	t.Run("nutanix", func(t *testing.T) { assertDisabled(t, &NutanixConfig{}) })
}

func TestKnownJSONKeysIncludesEmbeddedFields(t *testing.T) {
	keys := knownJSONKeys(reflect.TypeOf(&BaremetalConfig{}))

	// Directly declared, and reached through RhelConfig -> EpelConfig and
	// RhelConfig -> RhsmConfig.
	for _, key := range []string{"disk_size", "rhel_username", "redhat_epel_rpm", "epel_rpm_gpg_key", "rhsm_org_id", "iso_url", "http_proxy", "extra_rpms"} {
		assert.Containsf(t, keys, key, "%s should be recognized", key)
	}

	assert.NotContains(t, keys, "not_a_real_key")
}

func TestWarnOnUnknownConfigKeys(t *testing.T) {
	tests := []struct {
		name       string
		in         string
		wantLogged []string
		wantQuiet  bool
	}{
		{
			name:       "unknown keys are reported",
			in:         `{"disk_size":"20480","typo_size":"1","another_bogus_key":""}`,
			wantLogged: []string{"another_bogus_key", "typo_size", "baremetal-config"},
		},
		{
			name:      "recognized keys stay quiet",
			in:        `{"disk_size":"20480","redhat_epel_rpm":"","rhel_username":"u","rhsm_org_id":"1"}`,
			wantQuiet: true,
		},
		{
			name:      "malformed json defers to the typed unmarshal error",
			in:        `{"disk_size":`,
			wantQuiet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logged bytes.Buffer
			log.SetOutput(&logged)
			defer log.SetOutput(os.Stderr)

			WarnOnUnknownConfigKeys([]byte(tt.in), &BaremetalConfig{}, "baremetal-config")

			if tt.wantQuiet {
				assert.Empty(t, logged.String())
				return
			}
			for _, want := range tt.wantLogged {
				assert.Contains(t, logged.String(), want)
			}
		})
	}
}
