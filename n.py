#!/usr/bin/env python

import pprint
import subprocess
import sys
import traceback
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
    """
    Program's entry point
    """
    if len(sys.argv) < 2:
        sys.exit(USAGE)

    command = sys.argv[1]
    opts = parse_opts_root(*CONFIG_DEFAULTS.items())

    # get config
    config_path = Path(opts['config-file']['value']).expanduser()
    config = read_config(config_path, opts)

    if command == 'help' or command == '-h' or command == '--help':
        print(USAGE)
    elif command == 'version' or command == '-V' or command == '-v' or command == '--version':
        print(f'{NAME} v{VERSION}')
    elif command == 'new':
        cmd_new(config, sys.argv[2:])
    else:
        sys.exit(f'ERROR: invalid command ‘{command}’. See the help command for more information.')

def cmd_new(config, args):
    """
    Create a new note entry
    """
    node_dir = get_new_node_dir(config['notes-dir'])

    try:
        Path(node_dir).mkdir(parents=True,exist_ok=True)
    except:
        traceback.print_exc()
        sys.exit()

    # parse args
    title = "No title provided..." # Default title
    for arg in args:
        if not arg.startswith('--'):
            title = arg
            break

    full_path = Path(node_dir, "README.org")
    with open(str(full_path), 'w') as f:
        f.write(f'#+title: {title}\n\n')

    subprocess.run([config['editor'], full_path])
    print(f'Created note at ‘{node_dir}’ with the title ‘{title}’')

def get_new_node_dir(notes_dir):
    i = 1
    result = Path(notes_dir, str(i))
    while result.is_dir():
        i += 1
        result = Path(notes_dir, str(i))

    return str(result)

def read_config(path, opts):
    """
    Read a YAML configuration file into a python dictionary.

    path: YAML file
    opts: user-provided options returned by ‘parse_opts_root()’
    """
    config = {}
    try:
        with open(path, 'r') as f:
            config = yaml.safe_load(f)
    except:
        config = CONFIG_DEFAULTS

    overrides = {}
    for key, value in enumerate(opts):
        if opts[value]['override']:
            overrides[value] = opts[value]['value']

    config = {**CONFIG_DEFAULTS, **config, **overrides}
    config['notes-dir'] = Path(config['notes-dir']).expanduser()

    return config

def parse_opts_root(*args):
    """
    Parse command line options (arguments that start with ‘--’)
    """
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
