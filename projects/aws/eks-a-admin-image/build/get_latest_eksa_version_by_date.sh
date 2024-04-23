#!/usr/bin/env bash

set -x
set -o errexit
set -o nounset
set -o pipefail

RELEASES_MANIFEST=$(curl --silent -L $EKSA_RELEASE_MANIFEST_URL)
SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd -P)"
source "${SCRIPT_ROOT}/build/lib/common.sh"
DATE=$(build::find::gnu_variant_on_mac date)

newest_release_date=0
newest_release_version=""

while IFS=$'\t' read -r date version _; do
    # date is something like '2022-05-05 13:05:34.038243612 +0000 UTC'
    # I can't get `date` to parse this format, so deleting the extra timezone first
    parsed_date=$($DATE -d "$(echo $date | awk '{print $1,$2,$3}')" +"%s")

    if [ $parsed_date -gt $newest_release_date ];
    then
        newest_release_date=$parsed_date
        newest_release_version=$version
    fi  
done < <(echo "$RELEASES_MANIFEST" | yq e '.spec.releases[] | [.date, .version] | @tsv' -)

if [ -z "${newest_release_version}" ]; then
    echo "Not valid release found"
    exit 1
fi

echo -n "$newest_release_version"