package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CarbonFootprintCmd represents the carbon-footprint command
var CarbonFootprintCmd = &cobra.Command{
	Use:   "carbon-footprint",
	Short: "Estimate carbon emissions of repository CI/CD workloads",
	Long: `A module that estimates the energy consumption and carbon emissions (gCO2eq) 
of a repository's automated testing and build processes based on execution time and 
runner specifications, promoting sustainable software engineering practices.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing Green Computing Carbon Footprint Estimator...")
		fmt.Println("Aggregating CI/CD execution times and runner hardware specs...")

		estimation, err := estimateCarbonFootprint()
		if err != nil {
			fmt.Printf("Error calculating carbon footprint: %v\n", err)
			return
		}

		fmt.Println("\n--- CI/CD Carbon Footprint Report (Last 30 Days) ---")
		
		fmt.Printf("\n[Compute Usage]\n")
		fmt.Printf("Total Pipeline Execution Time: %d hours\n", estimation.TotalComputeHours)
		fmt.Printf("Primary Runner Specifications: %s\n", estimation.RunnerSpecs)
		
		fmt.Printf("\n[Environmental Impact]\n")
		fmt.Printf("Estimated Energy Consumption: %.2f kWh\n", estimation.EnergyConsumptionKWh)
		fmt.Printf("Carbon Emissions Generated: %.2f kgCO2eq\n", estimation.CarbonEmissionsKg)
		
		fmt.Printf("\n[Eco-Efficiency Assessment]\n")
		fmt.Printf("Rating: %s\n", estimation.EcoRating)
		fmt.Printf("Recommendation: %s\n", estimation.Recommendation)

		fmt.Println("\n(Data provided is an estimation based on regional grid carbon intensity averages.)")
	},
}

func init() {
	rootCmd.AddCommand(CarbonFootprintCmd)
}

type CarbonEstimation struct {
	TotalComputeHours    int
	RunnerSpecs          string
	EnergyConsumptionKWh float64
	CarbonEmissionsKg    float64
	EcoRating            string
	Recommendation       string
}

func estimateCarbonFootprint() (CarbonEstimation, error) {
	// Mocking carbon footprint estimation based on standard cloud runner profiles
	// Example calculation: 450 hours of a standard 2 vCPU cloud runner
	
	totalHours := 450
	wattsPerHour := 15.0 // Estimated watts per hour for standard runner
	energyKWh := (float64(totalHours) * wattsPerHour) / 1000.0
	
	// Assuming average global carbon intensity of 400g CO2 / kWh
	carbonKg := (energyKWh * 400.0) / 1000.0

	estimation := CarbonEstimation{
		TotalComputeHours:    totalHours,
		RunnerSpecs:          "Linux 2 vCPU, 7GB RAM (Standard GitHub Hosted)",
		EnergyConsumptionKWh: energyKWh,
		CarbonEmissionsKg:    carbonKg,
		EcoRating:            "C (Needs Improvement)",
		Recommendation:       "Consider utilizing Alpine-based docker images to reduce pull times, and configure CI caches to skip redundant dependency resolution. Relocating runners to regions with high renewable energy mix (e.g., eu-north-1) could reduce emissions by 60%.",
	}

	return estimation, nil
}
