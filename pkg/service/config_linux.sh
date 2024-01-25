#!/bin/sh

if [ -n "$GOFILE" ]; then
	exec > "$GOFILE~"
	trap "rm -f '$GOFILE~'" EXIT
fi

cat <<EOT
//go:build linux

package ${GOPACKAGE:-${PWD##*/}}

//go:generate $0
EOT

#
# FooScript
#
for x in Systemd Upstart Sysv OpenRC; do
	key="${x}Script"
	name="$(echo "$x" | tr 'A-Z' 'a-z')"

	cat <<EOT

// $key is the custom $name script.
func (cfg *Config) $key() string {
	return cfg.GetStringOption("$key", "")
}

// Set$key sets the custom $name script.
func (cfg *Config) Set$key(script string) {
	cfg.SetOption("$key", script)
}
EOT
done

if [ -n "$GOFILE" ]; then
	mv "$GOFILE~" "$GOFILE"
fi
