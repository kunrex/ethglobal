import dotenv from "dotenv"

import "hardhat-deploy"
import "hardhat-deploy-ethers"
import "@nomicfoundation/hardhat-toolbox"

dotenv.config()

const PRIVATE_KEY = process.env.PRIVATE_KEY
const CALIBRATION_URL = process.env.CALIBRATION_URL
const CHAIN_ID = parseInt(process.env.CHAIN_ID)

const config = {
    solidity: "0.8.17",
    defaultNetwork: "calibration",
    networks: {
        calibration: {
            chainId: CHAIN_ID,
            url: CALIBRATION_URL,
            accounts: [PRIVATE_KEY]
        },
    },
};

export default config;