// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

contract ProjectRegistry {
    struct Project {
        bytes cid;
        bytes metaData;
    }

    mapping (bytes32 => Project) private projects;

    function setProject(bytes32 index, bytes memory cid, bytes memory metaData) public {
        projects[index] = Project(cid, metaData);
    }

    function getProject(bytes32 index) public view returns (bytes memory, bytes memory, bool) {
        Project storage p = projects[index];
        return (p.cid, p.metaData, p.cid.length > 0);
    }

    function setMetaData(bytes32 index, bytes memory metaData) public {
        projects[index].metaData = metaData;
    }

    function getMetaData(bytes32 index) public view returns (bytes memory, bool) {
        bytes memory metaData = projects[index].metaData;
        return (metaData, metaData.length > 0);
    }
}