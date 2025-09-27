async function main() {
    const Contract = await ethers.getContractFactory("ProjectRegistry");
    const contract = await Contract.deploy();

    await contract.waitForDeployment();

    console.log("contract address:", await contract.getAddress());
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});