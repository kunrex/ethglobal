# On-Chain Remote Git Server
Just use git to permanently, securely, and seamlessly push major releases of your codebase on FileCoin, using Lighthouse API, with 0 downtime.
Framework to upload major version releases of projects to FileCoin through LightHouse API in a secure and seamless manner just using git
This ensures 0 downtime and acts as a decentralized cold storage for git projects

# Installation and Setup
clone the repository and run
1. Clone the project
2. Build the docker image `docker build -t dgit .`
3. Run the docker image `docker run --rm --name dgit dgit`
4. Use example.env to make `.env` and put your lighthouse api key there
5. Use docker inspect to find the ip of the container
6. Get the generated wallet using `ssh git@ip /ccg address` (Save this wallet!)
7. Send some funds to the wallet (for gas)

# Usage
1. Create a new repo on the remote using `ssh git@ip "mkdir repo && cd repo && git init"`
2. Add the remote to your local repo `git remote add storage git@ip:repo`
3. Now you can use simple git commands with this remote `git push --remote storage` `git clone git@ip:repo`

Even if the docker container dies, it is stateless except the wallet and api key, you can run it again and it will work without requiring any fixing!
