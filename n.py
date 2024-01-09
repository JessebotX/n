#!/usr/bin/env python

import pprint
import sys
import yaml

from pathlib import Path

NAME = sys.argv[0]
VERSION = '0.1.0'
USAGE = f'''
SYNOPSIS
--------
`{NAME} <COMMAND(required)> [GLOBAL_OPTIONS(optional)...]`

COMMANDS
--------
`help`, `-h`, `--help`
    - print usage information

`version`, `-V`, `--version`
    - print current version and configuration

`new [TITLE(optional)]`, `text [TITLE(optional)]`
    - Create a new text-based note entry.
    - Examples:
      ```
      {NAME} new
      ```
      Creates a new text note entry with default title.

      ```
      {NAME} new "Hello, world!"
      ```
      Creates a new text note entry with the provided title ‘Hello, world!’

`r <OBJECT_LINK(required)> [TITLE(optional)]`, `ref <OBJECT_LINK(required)> [TITLE(optional)]`
    - Create a note entry that acts as an entry for other text entries to reference to
    - OBJECT_LINK can be a local file or an actual link
      - In the case of a local file, the local file will be copied into the same directory as the note entry
    - Examples:
      ```
      {NAME} ref https://some-article-online.edu
      ```
      Creates a new reference note entry pointing to the provided link ‘https://some-article-online.edu’

      ```
      {NAME} ref ~/Pictures/picture.png "Reference: This is a picture"
      ```
      Creates a new reference note embedding a local file object (downloaded image, audio, video, etc.) with a provided title ‘Reference: This is a picture’

GLOBAL_OPTIONS
--------------
- `--config-file=<path(required)>`
  change location of configuration file

  Example: `{NAME} version --config-file=~/Documents/another-n-config.yml`
- ``

EXAMPLES
--------

'''
CONFIG_DEFAULTS = {
    'config-file': '~/.config/n/config.yml',
    'editor': 'nano',
    'notes-dir': '~/Documents/n-data'
}

def main():
    if len(sys.argv) < 2:
        sys.exit(USAGE)

    command = sys.argv[1]
    opts = parse_opts_root(*CONFIG_DEFAULTS.items())

    # get config
    config = CONFIG_DEFAULTS
    try:
        path = Path(opts['config-file']['value']).expanduser()
        config = read_config(path, opts)
    except FileNotFoundError:
        pass

    pprint.pprint(config) # tmp

    if command == 'help' or command == '-h' or command == '--help':
        print(USAGE)
    elif command == 'version' or command == '-V' or command == '-v' or command == '--version':
        print(f'{NAME} v{VERSION}')
    elif command == 'new':
        cmd_new(sys.argv[2:])
    else:
        sys.exit(f'ERROR: invalid command ‘{command}’. See the help command for more information.')

def cmd_new(args):
    """
    Create a new note entry
    """
    print('Command ‘new’')
    print('ARGS', args)

def read_config(path, opts):
    """
    Read a YAML configuration file into a python dictionary.

    path: YAML file
    opts: user-provided options returned by ‘parse_opts_root()’
    """
    config = {}
    with open(path, 'r') as f:
        config = yaml.safe_load(f)

    overrides = {}
    for key, value in enumerate(opts):
        if opts[value]['override']:
            overrides[value] = opts[value]['value']

    config = {**CONFIG_DEFAULTS, **config, **overrides}

    return config

def parse_opts_root(*args):
    result = {}
    for arg in args:
        result[arg[0]] = {
            'value': arg[1],
            'override': False
        }

    for key, value in enumerate(sys.argv):
        for arg in args:
            if value.startswith(f'--{arg[0]}'):
                result[arg[0]] = {
                    'value': value[len(f'--{arg[0]}='):],
                    'override': True
                } # count a ‘=’ character as well
                if len(result[arg[0]]) == 0:
                    sys.exit(f'ERROR: incorrect argument formatting for ‘--{arg[0]}’')


    return result

if __name__ == '__main__':
    main()
