#!/bin/sh

if [ -n "$GOFILE" ]; then

	exec > "$GOFILE~"
	trap "rm -f $GOFILE~" EXIT
fi

cat <<EOT
package $GOPACKAGE

//go:generate $0

import "github.com/spf13/cobra"
EOT

for x in \
	"NoArgs:returns an error if any args are included" \
	"OnlyValidArgs:returns an error if there are any positional argument that are not in [Config.ValidArgs]" \
	"ArbitraryArgs:never returns an error" \
	; do

	fn="${x%%:*}"

	echo
	echo "$x" | sed -e 's/:/ /' | fmt -w 70 | sed -e 's:^:// :'

	cat <<EOT
func $fn(cmd *cobra.Command, args []string) error {
	return cobra.$fn(cmd, args)
}
EOT
done

for x in \
	"MaximumNArgs:returns an error if there are more than N args" \
	"MinimumNArgs:returns an error if there are fewer than N args" \
	"ExactArgs:returns an error if there are not N args" \
	; do

	fn="${x%%:*}"

	echo
	echo "$x" | sed -e 's/:/ /' | fmt -w 70 | sed -e 's:^:// :'

cat <<EOT
func $fn(n int) cobra.PositionalArgs {
	return cobra.$fn(n)
}
EOT
done

if [ -n "$GOFILE" ]; then
	mv "$GOFILE~" "$GOFILE"
fi
