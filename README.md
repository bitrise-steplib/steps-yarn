# Run yarn command

[![Step changelog](https://shields.io/github/v/release/bitrise-community/steps-yarn?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-community/steps-yarn/releases)

Runs `yarn` with the given command and args.

<details>
<summary>Description</summary>


Yarn is a package manager that is compatible with the npm registry. Download your app's dependencies via yarn by using this Step.

This Step uses Corepack to automatically manage the Yarn version. The Yarn version is determined by the `packageManager` field in your `package.json` file or the `.yarnrc`/`.yarnrc.yml` configuration file.

### Configuring the Step

To use the Step, you need to configure your dependencies in your `package.json` file.

1. (Optional) Specify the Yarn version in your `package.json`:
   ```json
   {
     "packageManager": "yarn@4.1.0"
   }
   ```
   If not specified, Corepack will use its default Yarn version.

2. Set a command in **The yarn command to run** input.

   If you leave the input blank, the Step will simply install your dependencies. You can find the other available command in [yarn's documentation](https://yarnpkg.com/lang/en/docs/cli/).

3. Set the arguments in the **Arguments for running yarn commands** input.

   You can specify multiple arguments. Check out the available arguments for each command in yarn's documentation.

### Troubleshooting

If the Step fails, run it again with verbose logging enabled. To do so, set the **Enable verbose logging** input to `yes`. Doing so allows yarn to output more information about the command you ran.

Make sure your commands and arguments are correct, and that your packages are correctly defined in the `package.json` file.

Note: This Step requires Corepack to manage Yarn versions. Corepack is bundled with Node.js from 14.19.0 and 16.9.0 until 24.x. For Node.js 25+, Corepack will be automatically installed via npm if not already present.

### Useful links

[Corepack documentation](https://github.com/nodejs/corepack#readme)
[Getting started with React Native apps](https://devcenter.bitrise.io/getting-started/getting-started-with-react-native-apps/)
[Running Detox tests on Bitrise](https://devcenter.bitrise.io/testing/running-detox-tests-on-bitrise/)

### Related Steps

[Run Cocoapods install](https://www.bitrise.io/integrations/steps/cocoapods-install)
[Run npm command](https://www.bitrise.io/integrations/steps/npm)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `workdir` | Working directory of the step. You can leave it empty to not change it.  |  | `$BITRISE_SOURCE_DIR` |
| `command` | Specify the command to run with `yarn`. For example `add`. Leave it blank to install dependencies.  |  |  |
| `args` | Arguments are added to the `yarn` command. You can specify multiple arguments, separated by a space character. For example `react` or `-dev` |  |  |
| `verbose_log` | Choose if debug logging is enabled.  | required | `no` |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-community/steps-yarn/pulls) and [issues](https://github.com/bitrise-community/steps-yarn/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
