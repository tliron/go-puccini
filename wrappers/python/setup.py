#!/usr/bin/env python3

# Supported environment variables:
#
# PUCCINI_GO_VERSION: set this to override the Go distribution version used to compile Puccini

import os, pathlib, setuptools, subprocess, tempfile, shutil, sys

root = pathlib.Path(__file__).parents[0]

with open(root / 'puccini' / '__init__.py') as f:
  globals_ = {}
  exec(f.read(), globals_)
  version = globals_['__version__']

with open(root / 'description.md') as f:
  description = f.read()

go_version = os.environ.get('PUCCINI_GO_VERSION', '1.22.4')
root = str(root).replace('"', '\\"')

script = '''\
# Install Go
curl https://go.dev/dl/go{go_version}.linux-amd64.tar.gz --silent --location | tar -xz
GO=$PWD/go/bin/go

# Build libpuccini
cd "{root}/puccini/go-source/library"
"$GO" build \
    -buildmode=c-shared \
    -o=../../libpuccini.so \
    -ldflags " \
        -X 'github.com/tliron/kutil/version.GitVersion={version}'"
'''.format(root=root, version=version, go_version=go_version)

t = tempfile.mkdtemp()
try:
  subprocess.check_output(('bash',  '-o', 'pipefail', '-euxc', script), cwd=t)
except subprocess.CalledProcessError as e:
  print(e.output)
  sys.exit(e.returncode)
finally:
  shutil.rmtree(t)

class Distribution(setuptools.dist.Distribution):
  def has_ext_modules(_): # https://stackoverflow.com/a/62668026
    return True

setuptools.setup(
  name='puccini',
  version=version,
  description='Puccini',
  long_description=description,
  long_description_content_type='text/markdown',
  license='Apache License 2.0',
  author='Tal Liron',
  author_email='tal.liron@gmail.com',
  url='https://github.com/tliron/puccini',
  download_url='https://github.com/tliron/puccini/releases',
  classifiers=(
    'Development Status :: 4 - Beta',
    'Intended Audience :: Developers',
    'License :: OSI Approved :: Apache Software License'),

  packages=('puccini',),
  package_dir={'puccini': 'puccini'},
  package_data={'puccini': ['libpuccini.so']}, # for bdist
  install_requires=('ard',),

  distclass=Distribution)
