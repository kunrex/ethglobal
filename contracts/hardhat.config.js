import dotenv from "dotenv"

import "hardhat-deploy"
import "hardhat-deploy-ethers"
import "@nomicfoundation/hardhat-toolbox"

dotenv.config()

const PRIVATE_KEY = process.env.PRIVATE_KEY

const config = {
    solidity: "0.8.17",
    defaultNetwork: "calibration",
    networks: {
        calibration: {
            chainId: 314159,
            url: "https://api.calibration.node.glif.io/rpc/v1",
            accounts: [PRIVATE_KEY]
        },
    },
};

export default config;