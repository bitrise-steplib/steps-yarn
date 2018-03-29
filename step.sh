#!/bin/bash

set -e

# Install Yarn if we are running Ubuntu
if [ -f /etc/lsb-release ]; then
  if which yarn >/dev/null; then
    echo "Yarn already installed."
  else
    echo "Yarn not installed. Installing..."
    curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list
    sudo apt-get update && sudo apt-get install yarn
  fi
fi

# Change the working dir if necessary
if [ ! -z "${workdir}" ] ; then
  echo "==> Switching to working directory: ${workdir}"
  cd "${workdir}"
  if [ $? -ne 0 ] ; then
    echo " [!] Failed to switch to working directory: ${workdir}"
    exit 1
  fi
fi

echo "Yarn version:"
yarn --version
echo ""

set -x
yarn ${command} ${args}
