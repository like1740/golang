pragma solidity ^0.8.0;

contract Voting {
   mapping(string => uint256) private votes;

   string[] private voteCandidates;

   mapping(string => bool) private candidateRecorded;

   function vote(string memory condidate) public {
        votes[condidate] += 1;

        if (!(candidateRecorded[condidate])) {
            voteCandidates.push(condidate);
            candidateRecorded[condidate] = true;
        }
   }

   function getVotes(string memory condidate) public view returns (uint256) {
      return votes[condidate];
   }

   function resetVotes() public {
        for (uint256 i = 0; i < voteCandidates.length; i++) {
            votes[voteCandidates[i]] = 0;
            candidateRecorded[voteCandidates[i]] = false;
        }
        delete voteCandidates;
   }
}