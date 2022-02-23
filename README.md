# Run yarn command

[![Step changelog](https://shields.io/github/v/release/bitrise-community/steps-yarn?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-community/steps-yarn/releases)

Runs `yarn` with the given command and args.

<details>
<summary>Description</summary>


Yarn is a package manager that is compatible with the npm registry. Download your app's dependencies via yarn by using this Step.

### Configuring the Step

To use the Step, you need to configure your dependencies in your `package.json` file.

1. Set a command in **The yarn command to run** input.

   If you leave the input blank, the Step will simply install your dependencies. You can find the other available command in [yarn's documentation](https://yarnpkg.com/lang/en/docs/cli/).

1. Set the arguments in the **Arguments for running yarn commands** input.

   You can specify multiple arguments. Check out the available arguments for each command in yarn's documentation.

You can also cache the contents of the node_modules directory by setting the **Cache node_modules** input to `yes`.

### Troubleshooting

If the Step fails, run it again with verbose logging enabled. To do so, set the **Enable verbose logging** input to `yes`. Doing so allows yarn to output more information about the command you ran.

Make sure your commands and arguments are correct, and that your packages are correctly defined in the `package.json` file.

### Useful links

[Getting started with React Native apps](https://devcenter.bitrise.io/getting-started/getting-started-with-react-native-apps/)
[Running Detox tests on Bitrise](https://devcenter.bitrise.io/testing/running-detox-tests-on-bitrise/)

### Related Steps

[Run Cocoapods install](https://www.bitrise.io/integrations/steps/cocoapods-install)
[Run npm command](https://www.bitrise.io/integrations/steps/npm)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `workdir` | Working directory of the step. You can leave it empty to not change it.  |  | `$BITRISE_SOURCE_DIR` |
| `command` | Specify the command to run with `yarn`. For example `add`. Leave it blank to install dependencies.  |  |  |
| `args` | Arguments are added to the `yarn` command. You can specify multiple arguments, separated by a space character. For example `react` or `-dev` |  |  |
| `cache_local_deps` | Select if the contents of node_modules directory should be cached.  `yes`: Mark local dependencies to be cached. `no`: Do not use cache.  All node_modules folders (recursively) located under the working directory will be cached. | required | `no` |
| `verbose_log` | Choose if debug logging is enabled.  | required | `no` |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-community/steps-yarn/pulls) and [issues](https://github.com/bitrise-community/steps-yarn/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
