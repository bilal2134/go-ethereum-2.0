package cap

// orchestrator.go: Multi-dimensional consistency orchestrator
// Responsible for dynamic adjustment of consistency levels and network partition prediction.

import (
	"context"
	"log"
	"time"
)

// Orchestrator handles adaptive consistency and network partition prediction.
type Orchestrator struct {
	consistencyLevel int
	partitionPred    PartitionPredictor
}

// PartitionPredictor predicts network partition probability.
type PartitionPredictor interface {
	// Predict returns estimated partition probability (0.0-1.0).
	Predict(ctx context.Context) (float64, error)
}

// NewOrchestrator creates a new Orchestrator.
func NewOrchestrator(consistencyLevel int, partitionPred PartitionPredictor) *Orchestrator {
	return &Orchestrator{
		consistencyLevel: consistencyLevel,
		partitionPred:    partitionPred,
	}
}

// AdjustConsistency dynamically adjusts the consistency level based on partition probability.
func (o *Orchestrator) AdjustConsistency(ctx context.Context) {
	// Get partition probability
	prob, err := o.partitionPred.Predict(ctx)
	if err != nil {
		log.Printf("[Orchestrator] Partition prediction error: %v", err)
		return
	}
	if prob > 0.5 {
		log.Println("[Orchestrator] Partition risk high, using Strong Consistency")
		o.consistencyLevel = int(StrongConsistency)
	} else {
		log.Println("[Orchestrator] Partition risk low, using Eventual Consistency")
		o.consistencyLevel = int(EventualConsistency)
	}
}

// CurrentLevel returns the orchestrator's current consistency level.
func (o *Orchestrator) CurrentLevel() ConsistencyLevel {
	return ConsistencyLevel(o.consistencyLevel)
}

// PredictPartition triggers partition prediction and logs the result.
func (o *Orchestrator) PredictPartition(ctx context.Context) {
	_, _ = o.partitionPred.Predict(ctx)
}

// NetworkTelemetry holds network metrics for partition prediction.
type NetworkTelemetry struct {
	LatencyMs   float64
	PacketLoss  float64
	Throughput  float64
	LastUpdated time.Time
}

// SimplePartitionPredictor predicts network partition probability using telemetry.
type SimplePartitionPredictor struct {
	Telemetry NetworkTelemetry
}

// Predict estimates the probability of a network partition.
func (s *SimplePartitionPredictor) Predict(ctx context.Context) (float64, error) {
	prob := 0.0
	if s.Telemetry.LatencyMs > 200 {
		prob += 0.4
	}
	if s.Telemetry.PacketLoss > 0.05 {
		prob += 0.4
	}
	if s.Telemetry.Throughput < 1.0 {
		prob += 0.2
	}
	log.Printf("[PartitionPredictor] Partition probability: %.2f", prob)
	return prob, nil
}
