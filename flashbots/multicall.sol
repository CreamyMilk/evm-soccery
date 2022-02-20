// SPDX-License-Identifier: CC-BY-NC-SA-2.5

pragma solidity ^0.6.12;
pragma experimental ABIEncoderV2;

contract flashbotsmultiCall {
    function multiCall(address[] memory targets, bytes[] memory data, uint256[] memory values, uint256 coinbaseBribe) public payable {
        require(targets.length == data.length && data.length == values.length, "Length mismatch");
        for(uint i = 0; i < targets.length; i++) {
            (bool status,) = targets[i].call{value:values[i]}(data[i]);
            require(status, "call failed");
        }
        block.coinbase.transfer(coinbaseBribe);
    }
    function coinbaseTransfer() public payable {
        block.coinbase.transfer(msg.value);
    }
    receive() external payable {}
    fallback() external {}
}

