# oidc-exchange

## About this plugin
This plugin does an OIDC token exchange with the JFrog server. The resulting output is an access token that can be used
to access JFrog with the identity mapped in the OIDC provider configuration.

## Installation with JFrog CLI
Installing the latest version:

`$ jf plugin install oidc-exchange`

Installing a specific version:

`$ jf plugin install oidc-exchange@version`

Uninstalling a plugin

`$ jf plugin uninstall oidc-exchange`

## Usage
### Commands
* exchange
    - Flags:
        - shout: Makes output uppercase **[Default: false]**
    - Example:
    ```
  $ jf hello-frog hello world --shout
  
  NEW GREETING: HELLO WORLD
  ```

### Environment variables
* HELLO_FROG_GREET_PREFIX - Adds a prefix to every greet **[Default: New greeting: ]**

## Additional info
None.

## Release Notes
The release notes are available [here](RELEASE.md).
