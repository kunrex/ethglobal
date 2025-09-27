// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

contract ProjectRegistry {
    mapping (bytes32 => bytes) private projects;
    mapping (bytes32 => bytes) private projectMetaData;

    function setProject(bytes32 index, bytes memory cid) public {
        projects[index] = cid;
    }

    function getProject(bytes32 index) public view returns (bytes memory, bool) {
        bytes memory cid = projects[index];
        return (cid, cid.length > 0);
    }

     function setMetaData(bytes32 index, bytes memory cid) public {
        projectMetaData[index] = cid;
     }

     function getMetaData(bytes32 index) public view returns (bytes memory, bool) {
         bytes memory cid = projectMetaData[index];
         return (cid, cid.length > 0);
     }
}