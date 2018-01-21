When working with AWS there is often a need to manage credentials for multiple accounts. This can be done with SaaS solutions like LastPas or 1Password, but they are generally not very good at working with the AWS console login page. Also there is limited support for command line credentials management for accessing AWS through the cli tool.

AWS defaults to using a plain-text credentials file to store keys for programmatic access. These can be accessed from the command line using named profiles. Although easy, not very secure. 

awscreds was born out of the need to be able to manage multiple aws accounts for different companies in an easy and secure way. Providing an painfree way change between credentials in a secur way.

## Usage Example
