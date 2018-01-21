[![Build Status][circleci-badge]][circleci]
[![Release][release-badge]][release]
[![GoReportCard][report-badge]][report]
[![License][license-badge]][license]

[circleci-badge]: https://circleci.com/gh/TWinsnes/awscreds.svg?style=shield&circle-token=473ea84884d0f0afc3c0b259df55f7e97700b4fe
[circleci]: https://circleci.com/gh/TWinsnes/awscreds
[release-badge]: http://img.shields.io/github/release/TWinsnes/awscreds/all.svg?style=flat
[release]:https://github.com/TWinsnes/awscreds/releases
[report-badge]: https://goreportcard.com/badge/github.com/TWinsnes/awscreds
[report]: https://goreportcard.com/report/github.com/TWinsnes/awscreds
[license-badge]: https://img.shields.io/github/license/twinsnes/awscreds.svg?style=flat
[license]: https://github.com/twinsnes/awscreds/license
[homebrew-badge]: https://img.shields.io/homebrew/v/awscreds.svg
[homebrew]: https://homewbrew.com

When working with AWS there is often a need to manage credentials for multiple accounts. This can be done with SaaS solutions like LastPas or 1Password, but they are generally not very good at working with the AWS console login page. Also there is limited support for command line credentials management for accessing AWS through the cli tool.

AWS defaults to using a plain-text credentials file to store keys for programmatic access. These can be accessed from the command line using named profiles. Although easy, not very secure. 

awscreds was born out of the need to be able to manage multiple aws accounts for different companies in an easy and secure way. Providing an painfree way change between credentials in a secur way.

## Usage Example
![Console Example](console_example.gif)