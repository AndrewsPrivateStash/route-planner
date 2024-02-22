This project aims to replace the dependence on proprietary computational software (Mathematica) to cluster, route and image rapid deployment data.


* _Pre Conditions_: tss.exe and in.txt are in the same directory when run
* _Post Condition_: out.txt and out_route.png are created in directory tss is run

---

### Flags
```
-f     {"in.txt"}  define the input file name
-o     {"out.txt"} define outfile file name
-t     (true}      generate a tour image when done
-r     {0.8}       retention rate for simulated annealing process
-s     {0}         starting node (zero index) to rotate result to; default is the first node provided
-m     {"auto"}    select optimization method to use; default is dynamic method selection based on node-set
-a     {""}        provide an anchor to rotate the results to. Expects a string comma separated eg. -a="Lat,Lon"
-cls   {0}         generate this number of clusters and produce separate files and image output. Skip routing
-fmt   {true}      format output. formatting includes center headers and order column, false to pipe
-ctr   {false}     create and route centroids instead of locations using common labels
```

### Optimization Methods
* `exh`		exhaustive method, tries all possible permutations (scales by n! eg. 12! = 479001600), system processes about 500k/s
* `opt`		2-Opt method with simulated annealing. This is the best approach, but is slow above 1000 nodes
* `resOpt`	2-Opt method with simulated annealing and restricted swap search. Slow after about 3000 nodes
* `bigOpt`	nearest neighbor pass, then a single pass of 2-Opt without simulated annealing slow after 10000 nodes
* `nn`		nearest neighbor method. Fast for all reasonable node-sets but low quality
* `nnMul`	nearest neighbor with multi-start. Tries nearest neighbor for all starting nodes and chooses best
* `none`	skip optimization

---
	
### Sample Usage

`$ tss.exe`

run with all default flags

`$ tss.exe -f mydata.txt -o myOutputFile.txt -t=false -a 47.782816,-122.343771`

run with custom names for input and output files, skip tour image generation and rotate to the node nearest 47.782816,-122.343771

`$ tss.exe -t=false -s 3`

skip tour image generation and rotate to the 4th node

`$ tss.exe -s 99 -m nnMul`

run using the nearest neighbor multi-start method and rotate to the 100th node

`$ tss.exe -m none -t=false -a 47.782816,-122.343771`

run skipping optimization and no image output and anchor to 47.782816,-122.343771

`$ tss.exe -cls 10`

perform k-means clustering with 10 clusters and skip routing

`$ tss.exe -cls 5 -fmt=false`

perform k-means clustering with 5 clusters and no output formatting

`$ tss.exe -ctr=true -a 47.782816,-122.343771`

route using groups defined by labels (centroids) and rotate to anchor
