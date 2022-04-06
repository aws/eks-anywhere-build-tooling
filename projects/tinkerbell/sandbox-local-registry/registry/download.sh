#!/usr/bin/env bash
# This script handles downloading os images from s3 and stages them in nginx.

set -ex

main() {
	local images_file="$1"
	# this confusing IFS= and the || is to capture the last line of the file if there is no newline at the end
	while IFS= read -r img || [ -n "${img}" ]; do
		# file is expected to have src and dst images delimited by a space
		local src_img
		src_img="$(echo "${img}" | cut -d' ' -f1)"

		local dst_img
		dst_img="$(echo "${img}" | cut -d' ' -f2)"

		if [ ! -f "${dst_img}" ]; then
			wget "${src_img}" -O "${dst_img}"
		else
			echo "File ${dst_img} already exists!"
		fi

	done <"${images_file}"
}

main "$1"