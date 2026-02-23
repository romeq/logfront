# Logfront

Extendable program that reads logs and notifies you about authentication failures and suspicious activities
through different services. 

## Contributing

See [`CONTRIBUTING.md`](./CONTRIBUTING.md).

## Docs

> Config contains 2 primary keys: `sources` and `services`.
> 
> *Sources* are log sources, *services (or consumers)* are notification services (such as ntfy.sh).
> 
> Each source can have multiple *consumers*, which are services that will be notified about authentication failures.
> 
> Each source and service define their own configuration. Their configuration keys are defined in their respective documentation.
> 
> For SSH and FTP, in addition to `consumers` key you can define `systemd` and `logfile` keys.
> Systemd is used to determine whether the logs are being read from systemd journal or from a file. **Therefore `logfile` mustn't be defined when `systemd` is set to `true`.**
> 
> Each source must have at least one consumer, and **consumers must be defined as a list of strings under the key `consumers`**
> 

 Example config:
 ```yaml
 sources:
   ssh:
     logfile: /var/log/auth.log
     consumers:
       - ntfy_sh
   ftp:
     systemd: true
     consumers:
       - ntfy_sh
       - service
   some_source:
     that:
       - defines: its
         own: configuration
     consumers:
       - service
      
 
 services:
   service:
     with:
       some:
         - very: interesting
           config: setup
       that:
         has: stuff
         
   ntfy_sh:
     urls:
       - ntfy.sh/some_ntfy_topic
       - raspberrypi.local/some_other_ntfy_topic # allows using selfhosted instances

 ```

## To-Do
- [ ] Thorough documentation
- [ ] Different log sources
  - [x] <span style="color:green">Protocol for creating new sources</span>
  - [ ] <span style="color:orange">SSH</span>
  - [ ] <span style="color:orange">FTP</span>
  - [ ] <a style="color:orange" href="https://github.com/fail2ban/fail2ban">Fail2ban</a>
- [ ] Different notification services
  - [x] <span style="color:green">Protocol for integrating new services</span>
  - [ ] <span style="color:orange">ntfy.sh</span>
  - [ ] <span style="color:orange">Discord</span>
  - [ ] <span style="color:orange">Telegram</span>
- [ ] CI/CD
  - [x] <span style="color:green">Git pre-commit hooks</span>
  - [x] <span style="color:green">Build core on GitHub Actions</span>
  - [x] <span style="color:green">Run golangci-lint on Github Actions</span>
  - [ ] <span style="color:orange">Test coverage</span>
- [ ] Frontend (low priority)
