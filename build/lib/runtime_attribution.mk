# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

RUNTIME_ATTRIBUTION_GENERATOR=$(OUTPUT_DIR)/generate_runtime_attribution.sh

$(RUNTIME_ATTRIBUTION_GENERATOR): $(BUILD_LIB)/generate_runtime_attribution.sh
	@mkdir -p $(@D)
	cp $< $@

LICENSES_TARGETS_FOR_PREREQ+=$(RUNTIME_ATTRIBUTION_GENERATOR)

.PHONY: runtime-attribution
runtime-attribution: # Regenerate the public runtime source notices from the image's RPMs
runtime-attribution: IMAGE_NAME=$(REPO)
runtime-attribution: DOCKERFILE_FOLDER=./docker/linux
runtime-attribution: IMAGE_CONTEXT_DIR=.
runtime-attribution: override IMAGE_PLATFORMS=linux/$(BUILDER_PLATFORM_ARCH)
runtime-attribution: IMAGE_TARGET=runtime-attribution
runtime-attribution: IMAGE_OUTPUT_TYPE=local
runtime-attribution: IMAGE_OUTPUT=dest=$(OUTPUT_DIR)/runtime-attribution
runtime-attribution: IMAGE_IMPORT_CACHE=
runtime-attribution: ENSURE_PREREQ=$(if $(or $(filter true,$(BUILDCTL_AVAILABLE)),$(filter false,$(BUILDX_AVAILABLE))),ensure-buildkitd-host,)
runtime-attribution: export USE_BUILDX=$(if $(filter ensure-buildkitd-host,$(ENSURE_PREREQ)),false,true)
runtime-attribution: $(RUNTIME_ATTRIBUTION_GENERATOR) | $$(ENSURE_PREREQ)
	@$(BUILDCTL)
	cp $(OUTPUT_DIR)/runtime-attribution/RUNTIME_ATTRIBUTION.txt RUNTIME_ATTRIBUTION.txt

# The existing attribution bot collects *ATTRIBUTION.txt outside _output.
attribution: runtime-attribution
