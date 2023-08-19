# osrsSSBingo Genetic Algorithm Notes
## Gather Data: 
Collect information about each participant, including their weighted score which we can manually decide on or write a program and manually validate and preferences for the two categories (skiller/pvmer).

## Define the Weights: 
Assign appropriate weights to the score (<total level> * total_level_weight) + (<pvm skill> * pvm_skill_weight) + bank value *<bank value weight> and the two preference categories based on their choices of PVM or skilling or both.

## Normalize the Data: 
Normalize the weighted scores and preference values to ensure fair comparison.

## Implement Team Formation Algorithm: 
Use an algorithm like a Genetic Algorithm to form teams based on the normalized data while considering the preferences of each player
## Calculate Team Score: 
Calculate the team's overall score based on the weighted score and preferences of the teams' members.

## Optimize: 
Iterate through the algorithm several times with different random initializations to find the optimal team distribution.
## Evaluate Balance: 
Check the balance of the teams to ensure they are reasonably equal in terms of weighted scores and preferences.
## Output: 
Present the final teams to the participants in a clear and understandable format.