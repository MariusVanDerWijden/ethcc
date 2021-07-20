// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SmallContract {
    string constant public symbol = "TEST";
    mapping(address => uint256) balances;
    event Event(address indexed from, address indexed to, uint256 tokens);
    
    constructor(uint256 supply) {
        balances[msg.sender] = supply;
        emit Event(address(0), msg.sender, supply);
    }

    function balance(address owner) public view returns (uint256 balance) {
        return balances[owner];
    }

    function transfer(address to, uint256 tokens) public returns (bool success)
    {
        require(balances[msg.sender] >= tokens, "insufficient balance");
        balances[msg.sender] -= tokens;
        balances[to] += tokens;
        emit Event(msg.sender, to, tokens);
        return balances[msg.sender] == 0;
    }
}