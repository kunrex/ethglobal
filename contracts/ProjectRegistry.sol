pragma solidity ^0.8.0;

contract ProjectRegistry {
    mapping(address => bytes32) private projects;

    function setProject(bytes32 cid) public {
        projects[msg.sender] = cid;
    }

    function getProject(address wallet) public view returns (bytes32) {
        return projects[wallet];
    }
}