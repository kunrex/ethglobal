async function main() {
    const Contract = await ethers.getContractFactory("MyContract");
    const contract = await Contract.deploy();
    await contract.deployed();

    console.log("contract address:", contract.address);
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});