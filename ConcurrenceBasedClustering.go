package ConcurrenceBasedClustering

// =============================================================================
// Basic Concepts:
// 	This package is used for clustering of nodes based on concurrences of pairs
//	of nodes. In this package, there are quality models for evaluation of such
//	clusterings, as well cluster initializers and cluster optimizers.
// References:
//	[Shared Near Neighbors] Jarvis, R. A., & Patrick, E. A. (1973). Clustering
//		using a similarity measure based on shared near neighbors. IEEE
//		Transactions on computers, 100(11), 1025-1034.
//	[DBSCAN] Ester, M., Kriegel, H. P., Sander, J., & Xu, X. (1996, August). A
//		density-based algorithm for discovering clusters in large spatial
//		databases with noise. In Kdd (Vol. 96, No. 34, pp. 226-231).
//	[ROCK] Guha, S., Rastogi, R., & Shim, K. (2000). ROCK: A robust clustering
//		algorithm for categorical attributes. Information systems, 25(5), 345-
//		366.
//	[Centric Local Outliers] Yu, J. X., Qian, W., Lu, H., & Zhou, A. (2006).
//		Finding centric local outliers in categorical/numerical spaces.
//		Knowledge and Information Systems, 9(3), 309-338.
//	[Louvain Algorithm & Modularity] Blondel, V. D., Guillaume, J. L., Lambiotte
//		, R., & Lefebvre, E. (2008). Fast unfolding of communities in large
//		networks. Journal of statistical mechanics: theory and experiment,
//		2008(10), P10008.
//	[Constant Potts Model] Traag, V. A., Van Dooren, P., & Nesterov, Y. (2011).
//		Narrow scope for resolution-limit-free community detection. Physical
//		Review E, 84(1), 016114.
//	[Leiden Algorithm] Traag, V. A., Waltman, L., & Van Eck, N. J. (2019). From
//		Louvain to Leiden: guaranteeing well-connected communities. Scientific
//		reports, 9(1), 1-12.
//	[Label Propagation Algorithm] Zhu, X., & Ghahramani, Z. (2002). Learning
//		from labeled and unlabeled data with label propagation.
//	[Girvan Newman Algorithm] Girvan, M., & Newman, M. E. (2002). Community
//		structure in social and biological networks. Proceedings of the national
//		academy of sciences, 99(12), 7821-7826.
//	[Clique Percolation Method] Palla, G., Derényi, I., Farkas, I., & Vicsek, T.
//		(2005). Uncovering the overlapping community structure of complex
//		networks in nature and society. nature, 435(7043), 814-818.
//	[Advanced Clique Percolation Method] Salatino, A. A., Osborne, F., & Motta,
//		E. (2018, May). AUGUR: forecasting the emergence of new research topics.
//		In Proceedings of the 18th ACM/IEEE on Joint Conference on Digital
//		Libraries (pp. 303-312).
//	[Sequential Clique Percolation Method] Kumpula, J. M., Kivelä, M., Kaski, K.
//		, & Saramäki, J. (2008). Sequential algorithm for fast clique
//		percolation. Physical review E, 78(2), 026109.
//	[SLINK] Sibson, R. (1973). SLINK: an optimally efficient algorithm for the
//		single-link cluster method. The computer journal, 16(1), 30-34.
//	[CLINK] Defays, D. (1977). An efficient algorithm for a complete link method
//		. The Computer Journal, 20(4), 364-366.
// =============================================================================

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// =============================================================================
// func init
// brief description: init the package
func init() {
	rand.Seed(time.Now().UnixNano())
}

// =============================================================================
// struct ConcurrenceModel
// brief description: This is a struct for concurrence models
type ConcurrenceModel struct {
	n                 uint
	sumConcurrences   uint
	sumConcurrencesOf []uint
	concurrences      map[uint]map[uint]uint
}

// =============================================================================
// func NewConcurrenceModel
// brief description: create a new ConcurrenceModel object
func NewConcurrenceModel() ConcurrenceModel {
	return ConcurrenceModel{
		n:                 0,
		sumConcurrences:   0,
		sumConcurrencesOf: []uint{},
		concurrences:      map[uint]map[uint]uint{},
	}
}

// =============================================================================
// func verifyConcurrences
// brief description: check whether the concurrences are valid.
// input:
//	n: the number of nodes
//	concurrences: a matrix that its element (i,j) is the frequency of the
//		concurrence of node i and node j. If no such concurrence exists, then
//		the	element is 0.
// output:
//	nothing, but will raise fatal exceptions otherwise.
func verifyConcurrences(n uint, concurrences map[uint]map[uint]uint) {
	maxNodeID := uint(0)
	for u, weightsOfU := range concurrences {
		if u > maxNodeID {
			maxNodeID = u
		}
		for v, weightUV := range weightsOfU {
			if v > maxNodeID {
				maxNodeID = v
			}
			weightsOfV, exists := concurrences[v]
			if !exists {
				log.Fatalln("Asymmetric concurrence")
			}
			weightVU, exists := weightsOfV[u]
			if !exists || weightVU != weightUV {
				log.Fatalln("Asymmetric concurrence")
			}
		}
	}
	if maxNodeID >= n {
		log.Fatalln("maxNodeID >= n")
	}
}

// =============================================================================
// func getSumConcurrencesOf
// brief description: Compute a vector that the i-th component of the vector is
//	the sum of concurrences connected to node i.
// input:
//	n: the number of nodes
//	concurrences: a matrix that its element (i,j) is the frequency of the
//		concurrence between node i and node j. If no such concurrence exists,
//		then the element is 0.
// output:
//	the vector mentioned in brief description.
func getSumConcurrencesOf(n uint, concurrences map[uint]map[uint]uint) []uint {
	// -------------------------------------------------------------------------
	// step 1:
	sumConcurrencesOf := make([]uint, n)
	for u := uint(0); u < n; u++ {
		mySum := uint(0)
		weightsOfU, exists := concurrences[u]
		if exists {
			for _, weightUV := range weightsOfU {
				mySum += weightUV
			}
		}
		sumConcurrencesOf[u] = mySum
	}

	// -------------------------------------------------------------------------
	// step 2: return the result
	return sumConcurrencesOf
}

// =============================================================================
// func (cm *ConcurrenceModel) SetConcurrences
// brief description: set the concurrences of cm
// input:
//	n: the number of nodes
//	concurrence: a matrix that its element (i,j) is the frequency of the concurrence
//		between node i and node j. If no such concurrence exists, then the
//		element is 0.
// output:
//	nothing.
func (cm *ConcurrenceModel) SetConcurrences(n uint,
	concurrence map[uint]map[uint]uint) {
	// -------------------------------------------------------------------------
	// step 1: check whether the concurrences are valid.
	verifyConcurrences(n, concurrence)

	// -------------------------------------------------------------------------
	// step 2: get the nodewise sum of weights
	sumConcurrencesOf := getSumConcurrencesOf(n, concurrence)

	// -------------------------------------------------------------------------
	// step 3: compute the sum of all weights
	sumConcurrences := uint(0)
	for _, value := range sumConcurrencesOf {
		sumConcurrences += value
	}

	// -------------------------------------------------------------------------
	// step 4: set the fields
	cm.n = n
	cm.sumConcurrences = sumConcurrences
	cm.sumConcurrencesOf = sumConcurrencesOf
	cm.concurrences = concurrence
}

// =============================================================================
// func (cm ConcurrenceModel) GetN
func (cm ConcurrenceModel) GetN() uint {
	return cm.n
}

// =============================================================================
// func (cm ConcurrenceModel) GetConcurrencesOf
// brief description: get the concurrences related to a node
// input:
//	i: a point ID
// output:
//	the frequency of the concurrence of i if exists, 0 otherwise
func (cm ConcurrenceModel) GetConcurrencesOf(i uint) map[uint]uint {
	weightsOfI, exists := cm.concurrences[i]
	if exists {
		return weightsOfI
	} else {
		return map[uint]uint{}
	}
}

// =============================================================================
// func (cm ConcurrenceModel) GetConcurrence
// brief description: get concurrence between i and j
// input:
//	i, j: two point IDs
// output:
//	the frequency of the concurrence between i and j if the edge exists, 0
//	otherwise
func (cm ConcurrenceModel) GetConcurrence(i, j uint) uint {
	weightIJ, exists := cm.GetConcurrencesOf(i)[j]
	if exists {
		return weightIJ
	} else {
		return uint(0)
	}
}

// =============================================================================
// func (cm ConcurrenceModel) GetCompleteCommunties
// brief description: first copy input communites to the result, then add all
//	isolated points into the result as single point communities.
// input:
//	communities: a list of clusters.
// output:
//	the complete communities with isolated points added as single point
//	communities.
func (cm ConcurrenceModel) GetCompleteCommunities(communities []map[uint]bool,
) []map[uint]bool {
	// -------------------------------------------------------------------------
	// step 1: copy the communties into the result and mark the points in the
	// communities.
	result := []map[uint]bool{}
	pointMarkers := make([]bool, cm.n)
	for i := uint(0); i < cm.n; i++ {
		pointMarkers[i] = false
	}
	for _, community := range communities {
		myCommunity := map[uint]bool{}
		for point, _ := range community {
			if point >= cm.n {
				log.Fatal(fmt.Sprintf("point %d > n = %d", point, cm.n))
			}
			if pointMarkers[point] {
				log.Fatal(fmt.Sprintf("point %d is in multiple communities", point))
			}
			myCommunity[point] = true
			pointMarkers[point] = true
		}
		result = append(result, myCommunity)
	}

	// -------------------------------------------------------------------------
	// step 2: add isolated points into the result as single point communities
	for i := uint(0); i < cm.n; i++ {
		if !pointMarkers[i] {
			result = append(result, map[uint]bool{i: true})
		}
	}

	// -------------------------------------------------------------------------
	// step 3: return the result
	return result
}

// =============================================================================
// func (cm ConcurrenceModel) Aggregate
// brief description: aggregates concurrences according to communities
// input:
//	communities: a list of clusters.
// output:
//	the aggregated ConcurrenceModel
func (cm ConcurrenceModel) Aggregate(communities []map[uint]bool,
) ConcurrenceModel {
	// -------------------------------------------------------------------------
	// step 1: set newN and create an empty newConcurrences
	newN := uint(len(communities))
	newConcurrences := map[uint]map[uint]uint{}
	for i := uint(0); i < newN; i++ {
		newConcurrences[i] = map[uint]uint{}
	}

	// -------------------------------------------------------------------------
	// step 2: scans through the communities to fill newConcurrences
	for i1 := uint(0); i1+1 < newN; i1++ {
		c1 := communities[i1]
		for i2 := i1 + 1; i2 < newN; i2++ {
			c2 := communities[i2]
			weightI1I2 := uint(0)
			for pt1, _ := range c1 {
				weightsOfPt1, exists := cm.concurrences[pt1]
				if !exists {
					continue
				}
				for pt2, _ := range c2 {
					weightPt1Pt2, exists := weightsOfPt1[pt2]
					if exists {
						weightI1I2 += weightPt1Pt2
					}
				}
			}
			if weightI1I2 > uint(0) {
				newConcurrences[i1][i2] = weightI1I2
				newConcurrences[i2][i1] = weightI1I2
			}
		}
	}

	// -------------------------------------------------------------------------
	// step 3: create a new ConcurrenceModel using these data
	newCM := NewConcurrenceModel()
	newCM.SetConcurrences(newN, newConcurrences)

	// -------------------------------------------------------------------------
	// step 4: return the new ConcurrenceModel
	return newCM
}

// =============================================================================
// func (cm ConcurrenceModel) InduceSimilarities
// brief description: induce similarities from concurrences.
// input:
//	nothing
// output:
//	A similarity matrix induced from concurrences.
func (cm ConcurrenceModel) InduceSimilarities() map[uint]map[uint]float64 {
	simMat := map[uint]map[uint]float64{}
	for u := uint(0); u < cm.n; u++ {
		row := map[uint]float64{u: 1.0}
		cu := 0.5 / float64(cm.sumConcurrencesOf[u])
		weightsOfU := cm.GetConcurrencesOf(u)
		for v, weightUV := range weightsOfU {
			cv := 0.5 / float64(cm.sumConcurrencesOf[v])
			row[v] = float64(weightUV) * (cu + cv)
		}
		simMat[u] = row
	}
	return simMat
}

// =============================================================================
// func (cm ConcurrenceModel) connectsWell
// brief description: check whether the concurrence graph connects a node well
//	in a partition of communities.
// input:
//	u: a node ID
//	cu: the communityID of u
//	communities: a list of clusters
//	r: a threshold
// output:
//	true if it connects well, false otherwise
func (cm ConcurrenceModel) connectsWell(u, cu uint, communities []map[uint]bool,
	r float64) bool {
	c := communities[cu]
	weightsOfU := cm.GetConcurrencesOf(u)
	x := 0.0
	for v, _ := range c {
		if v == u {
			continue
		}
		weightUV, exists := weightsOfU[v]
		if exists {
			x += float64(weightUV)
		}
	}
	return x >= r*float64(len(c)-1)
}

// =============================================================================
// func (cm ConcurrenceModel) connectsWellBetween
// brief description: check whether the concurrence graph connects a node well
//	in a partition of communities.
// input:
//	cu, cv: two communityIDs
//	communities: a list of clusters
//	r: a threshold
// output:
//	true if it connects well, false otherwise
func (cm ConcurrenceModel) connectsWellTo(u, cu uint, communities []map[uint]bool,
	r float64) bool {
	c := communities[cu]
	weightsOfU := cm.GetConcurrencesOf(u)
	x := 0.0
	for v, _ := range c {
		if v == u {
			continue
		}
		weightUV, exists := weightsOfU[v]
		if exists {
			x += float64(weightUV)
		}
	}
	return x >= r*float64(len(c)-1)
}

// =============================================================================
// interface QualityModel
// brief description: This is an interface for quality models
type QualityModel interface {
	// The first two methods are parts of ConcurrenceModel. Therefore, for
	// those structs merged with ConcurreneModel, they already have these two
	// methods
	GetN() uint
	GetCompleteCommunities(communities []map[uint]bool) []map[uint]bool

	// This method is simiar to that of ConcurrenceModel. The difference is the
	// return value.
	Aggregate(communities []map[uint]bool) QualityModel

	// The last two methods are new to QualityModel. The implementations of this
	// interface must implement them.
	Quality(communities []map[uint]bool) float64
	DeltaQuality(communities []map[uint]bool, u, oldCu, newCu uint) float64
}

// =============================================================================
// struct Modularity
// brief introduction: this is an implementation of the famous Modularity
// 	quality model for network clustering
type Modularity struct {
	r float64
	ConcurrenceModel
}

// =============================================================================
// func NewModularity
// brief description: create a new Modularity
// input:
//	r: a threshold of modularity
func NewModularity(r float64) Modularity {
	return Modularity{
		r:                r,
		ConcurrenceModel: NewConcurrenceModel(),
	}
}

// =============================================================================
// func (qm *Modularity) Aggregate
func (qm Modularity) Aggregate(communities []map[uint]bool) QualityModel {
	return QualityModel(Modularity{qm.r, qm.ConcurrenceModel.Aggregate(communities)})
}

// =============================================================================
// func (qm *Modularity) Quality
// brief description: this implements Quality for interface QualityModel
// input:
//	communities: a list of clusters.
// output:
//	the value of Modularity
func (qm Modularity) Quality(communities []map[uint]bool) float64 {
	// -------------------------------------------------------------------------
	// step 1: compute 1/m and r/m
	oneOverM := 1.0 / float64(qm.sumConcurrences)
	rOverM := qm.r * oneOverM

	// -------------------------------------------------------------------------
	// step 2: compute modularity using the following equation:
	// modularity = 1/m sum_{i,j} (w_{i,j} - k_i * k_j * r/m) delta(c_i, c_j),
	// where:
	//	1/m = oneOverM,
	//	w_{i,j} = concurrence[i][j],
	//	k_u = nodewiseSumWeights[u],
	//	delta(s,t) = 0 if s != t, 1 if s == t.
	//	c_u = the community ID of u, i.e., communities[c][u] == true
	result := 0.0
	for _, c := range communities {
		for i, _ := range c {
			ki := float64(qm.sumConcurrencesOf[i])
			for j, _ := range c {
				kj := float64(qm.sumConcurrencesOf[j])
				result += float64(qm.GetConcurrence(i, j)) - rOverM*ki*kj
			}
		}
	}
	result *= oneOverM

	// -------------------------------------------------------------------------
	// step 3: return the result
	return result
}

// =============================================================================
// func (qm *Modularity) DeltaQuality
// brief description: this implements DeltaQuality for interface QualityModel
// input:
//	communities: a list of clusters.
//	u: a node ID, 0 <= u < n.
//	oldCu: the ID of the cluster u currently locates in.
//	newCu: the ID of the cluster u wants to move in.
// output:
//	The change amount of modularity.
// output:
//	the value of Modularity
func (qm Modularity) DeltaQuality(communities []map[uint]bool,
	u, oldCu, newCu uint) float64 {
	// -------------------------------------------------------------------------
	// step 1: check whether oldCu and newCu are the same one.
	// no change if oldCu == newCu
	if oldCu == newCu {
		return 0.0
	}

	// -------------------------------------------------------------------------
	// step 2: compute 1/m and r/m
	oneOverM := 1.0 / float64(qm.sumConcurrences)
	rOverM := qm.r * oneOverM

	// -------------------------------------------------------------------------
	// step 3: compute delta modularity. Note that:
	// modularity = 1/m sum_{i,j} (w_{i,j} - k_i * k_j * 1/m) delta(c_i, c_j),
	// where:
	//	1/m = oneOverM,
	//	w_{i,j} = concurrence[i][j],
	//	k_u = nodewiseSumWeights[u],
	//	delta(s,t) = 0 if s != t, 1 if s == t.
	//	c_u = the community ID of u, i.e., communities[c][u] == true
	// therfore:
	// delta modularity =
	//	1/m sum_{j in community newCu} (w_{u,j} - k_u * k_j * r/m)
	//	- 1/m sum_{j in community oldCu, j != i} (w_{u,j} - k_u * k_j * r/m)
	// (3.1) fetch weights of u and k_u
	weightsOfU := qm.GetConcurrencesOf(u)
	ku := float64(qm.sumConcurrencesOf[u])

	// (3.2) add to result the change at the new community of u
	result := 0.0
	newCommunityOfU := communities[newCu]
	for j := range newCommunityOfU {
		weightUJ, exists := weightsOfU[j]
		if !exists {
			weightUJ = uint(0)
		}
		kj := float64(qm.sumConcurrencesOf[j])
		result += float64(weightUJ) - rOverM*ku*kj
	}

	// (3.3) subtract from result the change at the old community of u
	oldCommunityOfU := communities[oldCu]
	for j := range oldCommunityOfU {
		if j == u {
			continue
		}
		weightUJ, exists := weightsOfU[j]
		if !exists {
			weightUJ = uint(0)
		}
		kj := float64(qm.sumConcurrencesOf[j])
		result -= float64(weightUJ) - rOverM*ku*kj
	}
	result *= oneOverM

	// -------------------------------------------------------------------------
	// step 4: return the result
	return result
}

// =============================================================================
// struct CPM
// brief introduction: this is an implementation of the famous Constant Potts
// 	quality model for network clustering
type CPM struct {
	r float64
	ConcurrenceModel
}

// =============================================================================
// func NewCPM
// brief description: create a new CPM
// input:
//	r: a threshold of CPM
func NewCPM(r float64) CPM {
	return CPM{
		r:                r,
		ConcurrenceModel: NewConcurrenceModel(),
	}
}

// =============================================================================
// func (qm CPM) Aggregate
func (qm CPM) Aggregate(communities []map[uint]bool) QualityModel {
	return QualityModel(CPM{qm.r, qm.ConcurrenceModel.Aggregate(communities)})
}

// =============================================================================
// func (qm *CPM) Quality
// brief description: this implements Quality for interface QualityModel
// input:
//	communities: a list of clusters.
// output:
//	the value of Modularity
func (qm CPM) Quality(communities []map[uint]bool) float64 {
	// -------------------------------------------------------------------------
	// step 1: compute CPM using the following equation:
	// CPM = sum_c (w_c - r size_c^2),
	// where:
	//	c is a community,
	//	size_c is the number of nodes in c,
	//	w_c is the sum of weight(i,j) for all i, j in c.
	result := 0.0
	for _, c := range communities {
		sizeC := float64(len(c))

		sumWeightsOfC := 0.0
		for i, _ := range c {
			weightsOfI := qm.GetConcurrencesOf(i)
			for j, _ := range c {
				weightIJ, exists := weightsOfI[j]
				if exists {
					sumWeightsOfC += float64(weightIJ)
				}
			}
		}

		result += sumWeightsOfC - qm.r*sizeC*sizeC
	}

	// -------------------------------------------------------------------------
	// step 3: return the result
	return result
}

// =============================================================================
// func (qm *CPM) DeltaQuality
// brief description: this implements DeltaQuality for interface QualityModel
// input:
//	communities: a list of clusters.
//	u: a node ID, 0 <= u < n.
//	oldCu: the ID of the cluster u currently locates in.
//	newCu: the ID of the cluster u wants to move in.
// output:
//	The change amount of modularity.
// output:
//	the value of Modularity
func (qm CPM) DeltaQuality(communities []map[uint]bool,
	u, oldCu, newCu uint) float64 {
	// -------------------------------------------------------------------------
	// step 1: check whether oldCu and newCu are the same one.
	// no change if oldCu == newCu
	if oldCu == newCu {
		return 0.0
	}

	// -------------------------------------------------------------------------
	// step 2: compute delta CPM.
	// CPM = sum_c (w_c - r size_c^2),
	// where:
	//	c is a community,
	//	size_c is the number of nodes in c,
	//	w_c is the sum of weight(i,j) for all i, j in c.
	// Therefore:
	// delta CPM = delta w_oldCu + delta w_newCu
	//	- r ((size_oldCu-1)^2 - size_oldCu^2)
	//	- r ((size_newCu+1)^2 - size_newCu^2)
	//	= delta w_oldCu + delta w_newCu - r (-2 size_oldCu + 1)
	//	- r (2 size_newCu + 1)
	//	= delta w_oldCu + delta w_newCu - 2 r(size_newCu - size_oldCu + 1)

	// (2.1) fetch weights of u
	weightsOfU := qm.GetConcurrencesOf(u)

	// (2.2) compute delta w_oldCu
	deltaWOldCu := 0.0
	oldCommunityOfU := communities[oldCu]
	for j := range oldCommunityOfU {
		if j == u {
			continue
		}
		weightUJ, exists := weightsOfU[j]
		if exists {
			deltaWOldCu -= float64(weightUJ)
		}
	}

	// (2.3) compute delta w_newCu
	deltaWNewCu := 0.0
	newCommunityOfU := communities[newCu]
	for j := range newCommunityOfU {
		weightUJ, exists := weightsOfU[j]
		if exists {
			deltaWNewCu += float64(weightUJ)
		}
	}

	// (2.4) compute size_oldCu and size_newCu
	sizeOldCu := float64(len(oldCommunityOfU))
	sizeNewCu := float64(len(newCommunityOfU))

	// (2.5) compute the result
	result := deltaWOldCu + deltaWNewCu - 2*qm.r*(sizeNewCu-sizeOldCu+1)

	// -------------------------------------------------------------------------
	// step 3: return the result
	return result
}

// =============================================================================
// func getCorePoints
// brief description: This is part of an implementation to the famous DBScan
//	algorithm: looking for all core points.
// input:
//	simMat: the similarity matrix. It must be symmetric, all elements 0~1, and
//		the diagonal elements are all 1.
//	eps: the radius of neighborhood.
//	minPts: Only if the neighborhood of a point contains at least minPt points
//		(the center point of the neighborhood included), the neighborhood is
//		called dense. Only dense neighborhoods are connected to communities.
// output:
//	A map of core points to their neighborhood densities.
func getCorePoints(simMat map[uint]map[uint]float64, eps float64,
	minPts uint) map[uint]uint {
	// -------------------------------------------------------------------------
	// step 1: compute the density of all points' neighborhoods
	n := uint(len(simMat))
	densities := make([]uint, n)
	for pt := uint(0); pt < n; pt++ {
		rowPt, exists := simMat[pt]
		if !exists {
			log.Fatal("Invalid similarity matrix")
		}
		density := uint(0)
		for _, similarity := range rowPt {
			if similarity+eps >= 1.0 {
				density++
			}
		}
		densities[pt] = density
	}

	// -------------------------------------------------------------------------
	// step 2: generate a list of points with dense neighborhoods
	corePts := map[uint]uint{}
	for pt, density := range densities {
		if density >= minPts {
			corePts[uint(pt)] = density
		}
	}

	// -------------------------------------------------------------------------
	// step 3: return the result
	return corePts
}

// =============================================================================
// func getNeighbors
// brief description: This is part of an implementation to the famous DBScan
//	algorithm: generating a list of core members and another list of noncore
//	neighbors for each core points.
// input:
//	simMat: the similarity matrix. It must be symmetric, all elements 0~1, and
//		the diagonal elements are all 1.
//	eps: the radius of neighborhood.
//	minPts: Only if the neighborhood of a point contains at least minPt points
//		(the center point of the neighborhood included), the neighborhood is
//		called dense. Only dense neighborhoods are connected to communities.
//	corePts: a map of core points to their neighborhood densities.
// output:
//	output 1: a list of the core neighbors for each core point.
//	output 2: a list of the noncore neighbors for each core point.
func getNeighbors(simMat map[uint]map[uint]float64, eps float64,
	minPts uint, corePts map[uint]uint) (map[uint]map[uint]bool,
	map[uint]map[uint]bool) {
	coreNeighbors := map[uint]map[uint]bool{}
	noncoreNeighbors := map[uint]map[uint]bool{}
	for pt, _ := range corePts {
		// create the rows of the results
		coreRow := map[uint]bool{}
		coreNeighbors[pt] = coreRow
		noncoreRow := map[uint]bool{}
		noncoreNeighbors[pt] = noncoreRow

		// read the row of similarity matrix
		simRow, rowExists := simMat[pt]
		if !rowExists {
			log.Fatal("invalid similarity matrix")
		}

		// scan through the row we just read
		for neighbor, similarity := range simRow {
			// skip pt itself
			if neighbor == pt {
				continue
			}
			// find points that locate within pt's neighborhood
			if similarity+eps >= 1.0 {
				_, isCorePoint := corePts[pt]
				if isCorePoint {
					coreRow[neighbor] = true
				} else {
					noncoreRow[neighbor] = true
				}
			}
		}
	}
	return coreNeighbors, noncoreNeighbors
}

// =============================================================================
// func DBScan
// brief description: This is an implementation to the famous DBScan algorithm.
// input:
//	n: the number of nodes.
//	concurrence: a matrix that its element (i,j) is the frequency of the concurrence
//		between node i and node j. If no such concurrence exists, then the
//		element is 0.
//	eps: the radius of neighborhood.
//	minPts: Only if the neighborhood of a point contains at least minPt points
//		(the center point of the neighborhood included), the neighborhood is
//		called dense. Only dense neighborhoods are connected to communities.
// output:
//	A list of clusters.
func DBScan(cm ConcurrenceModel, eps float64, minPts uint) []map[uint]bool {
	// -------------------------------------------------------------------------
	// step 1: initialize auxiliary data structures
	communityIDs := map[uint]uint{}
	communities := []map[uint]bool{}

	// -------------------------------------------------------------------------
	// step 2: build the similarity matrix
	simMat := cm.InduceSimilarities()
	// -------------------------------------------------------------------------
	// step 3: find all core points and their neighborhood densities
	corePts := getCorePoints(simMat, eps, minPts)

	// -------------------------------------------------------------------------
	// step 4: find neighbors for each core point
	coreNeighbors, noncoreNeighbors := getNeighbors(simMat, eps, minPts, corePts)

	// -------------------------------------------------------------------------
	// step 5: loop until all core points are in communities
	n := cm.GetN()
	for {
		// (5.1) prepare an ID for the new community
		c := uint(len(communities))

		// (5.2) find the densist unassigned core point as the center point of
		// the new cluster
		centerPt := n
		centerDensity := uint(0)
		for pt, density := range corePts {
			// skip those points that have already been assigned into community
			_, exists := communityIDs[pt]
			if exists {
				continue
			}

			// check whether with the currently most dense neighborhood
			if density > centerDensity {
				centerPt = pt
				centerDensity = density
			}
		}

		// (5.3) stop the loop if not new centerPt is found
		if centerPt == n {
			break
		}

		// (5.4) officially create the community
		newCommunity := map[uint]bool{centerPt: true}
		communities = append(communities, newCommunity)
		communityIDs[centerPt] = c

		// (5.5) iteratively append neighbors to the new community
		boundary := map[uint]bool{centerPt: true}
		for len(boundary) > 0 {
			newBoundary := map[uint]bool{}
			for bpt, _ := range boundary {
				bptNoncoreNeighbors, exists := noncoreNeighbors[bpt]
				if exists {
					for neighbor, _ := range bptNoncoreNeighbors {
						newCommunity[neighbor] = true
						communityIDs[neighbor] = c
					}
				}
				bptCoreNeighbors, exists := coreNeighbors[bpt]
				if !exists {
					continue
				}
				for neighbor, _ := range bptCoreNeighbors {
					// skip those already in this community
					_, alreadyIn := newCommunity[neighbor]
					if !alreadyIn {
						newBoundary[neighbor] = true
						newCommunity[neighbor] = true
						communityIDs[neighbor] = c
					}
				}
			}
			boundary = newBoundary
		}
	}

	// -------------------------------------------------------------------------
	// step 6: return the result
	return communities
}

// =============================================================================
// func flattenCommunities
// brief description: expand the aggregated concurrence graph's communities at
//	the original concurrence graph.
// input:
//	aggCommunities: the aggregated concurrence graph's communities
//	communities: the original concurrence graph's communities
// output:
//	the flatten communities
func flattenCommunities(aggCommunities, communities []map[uint]bool,
) []map[uint]bool {
	result := []map[uint]bool{}
	for _, aggC := range aggCommunities {
		newC := map[uint]bool{}
		for idxC, _ := range aggC {
			c := communities[idxC]
			for pt, _ := range c {
				newC[pt] = true
			}
		}
		result = append(result, newC)
	}
	return result
}

// =============================================================================
// func Louvain
// brief description: Louvain algorithm for partition optimization of
//	concurrence graphs.
// input:
//	qm: a quality model.
//	communities: a list of clusters.
//	opts: an optional list of options
// output:
//	the optimized communities that maximizes quality
// note:
//	If the input communities is empty, this function will act as the classical
//	Louvain algorithm that uses single point communities as the initial
//	communities.
func Louvain(qm QualityModel, communities []map[uint]bool, opts ...string,
) []map[uint]bool {
	// step 1: parsing options
	useSeqSelector := true
	multiResolution := true
	shuffle := false
	for _, opt := range opts {
		switch opt {
		case "priority selector":
			useSeqSelector = false
		case "sequential selector":
			useSeqSelector = true
		case "single resolution":
			multiResolution = false
		case "multiple resolution":
			multiResolution = true
		case "shuffle":
			shuffle = true
		case "no shuffle":
			shuffle = false
		}
	}

	// -------------------------------------------------------------------------
	// step 2: complete communities with isolated points added as single point
	// communities.
	result := qm.GetCompleteCommunities(communities)
	n := qm.GetN()

	// -------------------------------------------------------------------------
	// step 3: get the community ID for each point
	communityIDs := make([]uint, n)
	for communityID, community := range result {
		for point, _ := range community {
			communityIDs[point] = uint(communityID)
		}
	}

	// -------------------------------------------------------------------------
	// step 4: iteratively scan through the points to find out what is the best
	// community for a point. If all points are in their best communities, stop
	// the iteration.
	m := uint(len(result))
	for {
		// (4.1) create the access order of points
		points := make([]uint, n)
		for i := 0; i < int(n); i++ {
			points[i] = uint(i)
		}

		// (4.2) optionally, shuffle the access order of points
		if shuffle {
			rand.Shuffle(int(n), func(i, j int) {
				points[i], points[j] = points[j], points[i]
			})
		}

		// (4.3) move points
		if useSeqSelector {
			done := true
			for _, u := range points {
				oldCu := communityIDs[u]
				bestDeltaQuality := 0.0
				bestNewCu := oldCu
				for newCu := uint(0); newCu < m; newCu++ {
					deltaQuality := qm.DeltaQuality(result, u, oldCu, newCu)
					if deltaQuality > bestDeltaQuality {
						bestDeltaQuality = deltaQuality
						bestNewCu = newCu
					}
				}

				if bestDeltaQuality > 0.0 {
					delete(result[oldCu], u)
					result[bestNewCu][u] = true
					communityIDs[u] = bestNewCu
					done = false
				}
			}
			if done {
				break
			}
		} else {
			bestDeltaQuality := 0.0
			bestU := uint(0)
			oldCBestU := communityIDs[0]
			bestNewCu := oldCBestU
			for _, u := range points {
				oldCu := communityIDs[u]
				for newCu := uint(0); newCu < m; newCu++ {
					deltaQuality := qm.DeltaQuality(result, u, oldCu, newCu)
					if deltaQuality > bestDeltaQuality {
						bestDeltaQuality = deltaQuality
						bestU = u
						oldCBestU = oldCu
						bestNewCu = newCu
					}
				}
			}
			if bestDeltaQuality == 0.0 {
				break
			}
			delete(result[oldCBestU], bestU)
			result[bestNewCu][bestU] = true
			communityIDs[bestU] = bestNewCu
		}
	}

	// -------------------------------------------------------------------------
	// step 5: remove empty communities
	oldResult := result
	result = []map[uint]bool{}
	for _, c := range oldResult {
		if len(c) > 0 {
			result = append(result, c)
		}
	}

	// -------------------------------------------------------------------------
	// step 6: if required, do the multi-resolution part
	if multiResolution {
		// ---------------------------------------------------------------------
		// (6.1) create aggregate network from the result
		newQM := qm.Aggregate(result)

		// ---------------------------------------------------------------------
		// (6.2) compute aggregated result from the aggregate network
		aggResult := Louvain(qm, result, opts...)

		// ---------------------------------------------------------------------
		// (6.3) check whether the new result has merged something. If it has,
		// then revise the result accordingly
		if uint(len(aggResult)) < newQM.GetN() {
			result = flattenCommunities(aggResult, result)
		}
	}

	// -------------------------------------------------------------------------
	// step 7: return the result
	return result
}

// =============================================================================
// func Leiden
// brief description: Leiden algorithm for partition optimization of
//	concurrence graphs.
// input:
//	qm: a quality model.
//	communities: a list of clusters.
//	opts: an optional list of options
// output:
//	the optimized communities that maximizes quality
// note:
//	If the input communities is empty, this function will act as the classical
//	Leiden algorithm that uses single point communities as the initial
//	communities.
func Leiden(qm QualityModel, communities []map[uint]bool, opts ...string,
) []map[uint]bool {
	// step 1: parsing options
	useSeqSelector := true
	multiResolution := true
	shuffle := false
	for _, opt := range opts {
		switch opt {
		case "priority selector":
			useSeqSelector = false
		case "sequential selector":
			useSeqSelector = true
		case "single resolution":
			multiResolution = false
		case "multiple resolution":
			multiResolution = true
		case "shuffle":
			shuffle = true
		case "no shuffle":
			shuffle = false
		}
	}

	// -------------------------------------------------------------------------
	// step 2: complete communities with isolated points added as single point
	// communities.
	result := qm.GetCompleteCommunities(communities)
	n := qm.GetN()

	// -------------------------------------------------------------------------
	// step 3: get the community ID for each point
	communityIDs := make([]uint, n)
	for communityID, community := range result {
		for point, _ := range community {
			communityIDs[point] = uint(communityID)
		}
	}

	// -------------------------------------------------------------------------
	// step 4: iteratively scan through the points to find out what is the best
	// community for a point. If all points are in their best communities, stop
	// the iteration.
	m := uint(len(result))
	for {
		// (4.1) create the access order of points
		points := make([]uint, n)
		for i := 0; i < int(n); i++ {
			points[i] = uint(i)
		}

		// (4.2) optionally, shuffle the access order of points
		if shuffle {
			rand.Shuffle(int(n), func(i, j int) {
				points[i], points[j] = points[j], points[i]
			})
		}

		// (4.3) move points
		if useSeqSelector {
			done := true
			for _, u := range points {
				oldCu := communityIDs[u]
				bestDeltaQuality := 0.0
				bestNewCu := oldCu
				for newCu := uint(0); newCu < m; newCu++ {
					deltaQuality := qm.DeltaQuality(result, u, oldCu, newCu)
					if deltaQuality > bestDeltaQuality {
						bestDeltaQuality = deltaQuality
						bestNewCu = newCu
					}
				}

				if bestDeltaQuality > 0.0 {
					delete(result[oldCu], u)
					result[bestNewCu][u] = true
					communityIDs[u] = bestNewCu
					done = false
				}
			}
			if done {
				break
			}
		} else {
			bestDeltaQuality := 0.0
			bestU := uint(0)
			oldCBestU := communityIDs[0]
			bestNewCu := oldCBestU
			for _, u := range points {
				oldCu := communityIDs[u]
				for newCu := uint(0); newCu < m; newCu++ {
					deltaQuality := qm.DeltaQuality(result, u, oldCu, newCu)
					if deltaQuality > bestDeltaQuality {
						bestDeltaQuality = deltaQuality
						bestU = u
						oldCBestU = oldCu
						bestNewCu = newCu
					}
				}
			}
			if bestDeltaQuality == 0.0 {
				break
			}
			delete(result[oldCBestU], bestU)
			result[bestNewCu][bestU] = true
			communityIDs[bestU] = bestNewCu
		}
	}

	// -------------------------------------------------------------------------
	// step 5: remove empty communities
	oldResult := result
	result = []map[uint]bool{}
	for _, c := range oldResult {
		if len(c) > 0 {
			result = append(result, c)
		}
	}

	// -------------------------------------------------------------------------
	// step 6: if required, do the multi-resolution part
	if multiResolution {
		// ---------------------------------------------------------------------
		// (6.1) create aggregate network from the result
		newQM := qm.Aggregate(result)

		// ---------------------------------------------------------------------
		// (6.2) compute aggregated result from the aggregate network
		aggResult := Leiden(qm, result, opts...)

		// -------------------------------------------------------------------------
		// (6.3) check whether the new result has merged something. If it has,
		// then revise the result accordingly
		if uint(len(aggResult)) < newQM.GetN() {
			result = flattenCommunities(aggResult, result)
		}
	}

	// -------------------------------------------------------------------------
	// step 7: return the result
	return result
}