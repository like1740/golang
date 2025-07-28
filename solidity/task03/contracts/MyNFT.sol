// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Context.sol";

contract MyNFT is ERC721, Ownable {
    uint256 private _tokenIdCounter;

    constructor() ERC721("MyNFT", "MNFT") {
        _tokenIdCounter = 0;
    }

    function mint(address to) public onlyOwner returns (uint256) {
        _tokenIdCounter++;
        uint256 tokenId = _tokenIdCounter;
        _safeMint(to, tokenId);
        return tokenId;
    }

    function burn(uint256 tokenId) public {
        require(isApprovedOrOwner(_msgSender(), tokenId),"ERC721: caller is not token owner or approved");
        _burn(tokenId);
    }

     function currentTokenId() public view returns (uint256) {
        return _tokenIdCounter;
    }
}