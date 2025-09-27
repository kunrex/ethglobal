pragma solidity ^0.8.17;

contract ProjectRegistry {
    mapping (bytes32 => bytes32) private projects;

    function setProject(bytes32 index, bytes32 cid) public {
        projects[index] = cid;
    }

    function getProject(bytes32 indedx) public view returns (bytes32, bool) {
        bytes32 cid = projects[indedx];
        if (cid == bytes32(0)) {
            return (cid, false);
        }

        return (cid, true);
    }
}