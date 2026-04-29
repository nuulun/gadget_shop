package seed

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"product-service/internal/model"
	"product-service/internal/repository"
)

// Run seeds the products table if it is empty.
// It uses the repository so it goes through the same data path as production code.
func Run(ctx context.Context, repo *repository.Repository) error {
	empty, err := repo.IsEmpty(ctx)
	if err != nil {
		return fmt.Errorf("seed: cannot check emptiness: %w", err)
	}
	if !empty {
		log.Println("[seed] products table already has data – skipping seed")
		return nil
	}

	log.Println("[seed] seeding products table …")
	products := buildProducts()
	if err := repo.BulkCreate(ctx, products); err != nil {
		return fmt.Errorf("seed: BulkCreate: %w", err)
	}
	log.Printf("[seed] inserted %d products ✓", len(products))
	return nil
}

// ─── Data ─────────────────────────────────────────────────────────────────────

type categorySpec struct {
	name   string
	brands []string
	images []string
	minP   float64
	maxP   float64
}

var specs = []categorySpec{
	{
		name:   "keyboard",
		brands: []string{"Logitech", "Razer", "HyperX", "Keychron", "Corsair"},
		images: []string{
			"https://images.unsplash.com/photo-1587829741301-dc798b83add3?w=400",
			"https://images.unsplash.com/photo-1595225476474-87563907a212?w=400",
			"https://images.unsplash.com/photo-1618384887929-16ec33fab9ef?w=400",
		},
		minP: 29, maxP: 250,
	},
	{
		name:   "mouse",
		brands: []string{"Logitech", "Razer", "SteelSeries", "Zowie", "Glorious"},
		images: []string{
			"https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?w=400",
			"https://images.unsplash.com/photo-1615663245857-ac93bb7c39e7?w=400",
			"https://images.unsplash.com/photo-1613141411244-0e4ac259d217?w=400",
		},
		minP: 19, maxP: 160,
	},
	{
		name:   "monitor",
		brands: []string{"Dell", "Asus", "LG", "Samsung", "BenQ"},
		images: []string{
			"https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?w=400",
			"https://images.unsplash.com/photo-1585792180666-f7347c490ee2?w=400",
			"https://images.unsplash.com/photo-1593640408182-31c70c8268f5?w=400",
		},
		minP: 149, maxP: 1200,
	},
	{
		name:   "laptop",
		brands: []string{"Asus", "Dell", "Apple", "Lenovo", "HP"},
		images: []string{
			"https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=400",
			"https://images.unsplash.com/photo-1525547719571-a2d4ac8945e2?w=400",
			"https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=400",
		},
		minP: 699, maxP: 3500,
	},
	{
		name:   "headset",
		brands: []string{"HyperX", "SteelSeries", "Razer", "Sony", "JBL"},
		images: []string{
			"https://images.unsplash.com/photo-1599669454699-248893623440?w=400",
			"https://images.unsplash.com/photo-1618366712010-f4ae9c647dcb?w=400",
			"https://images.unsplash.com/photo-1583394838336-acd977736f90?w=400",
		},
		minP: 29, maxP: 350,
	},
	{
		name:   "webcam",
		brands: []string{"Logitech", "Razer", "Microsoft", "Elgato", "OBSBOT"},
		images: []string{
			"https://images.unsplash.com/photo-1587826080692-f439cd0b70da?w=400",
			"https://images.unsplash.com/photo-1612831455359-970e23a1e4e9?w=400",
			"https://images.unsplash.com/photo-1593642632559-0c6d3fc62b89?w=400",
		},
		minP: 39, maxP: 300,
	},
}

var adjectives = []string{
	"Pro", "Elite", "Ultra", "Max", "Gaming", "Wireless", "RGB", "Mechanical",
	"Compact", "Silent", "Ergonomic", "Portable", "Premium",
}

var descriptions = map[string][]string{
	"keyboard": {
		"Tactile mechanical switches with satisfying click feedback for every keystroke.",
		"Full-size layout with programmable RGB backlighting and media controls.",
		"Compact tenkeyless design perfect for desktop setups with limited space.",
	},
	"mouse": {
		"High-precision optical sensor with adjustable DPI from 200 to 25,600.",
		"Ergonomic right-handed design with braided cable and onboard memory.",
		"Ultra-lightweight honeycomb shell for fatigue-free gaming sessions.",
	},
	"monitor": {
		"IPS panel with 1ms response time and 144Hz refresh rate for smooth gaming.",
		"4K UHD resolution with HDR400 support and factory-calibrated color accuracy.",
		"Ultrawide 21:9 curved display for immersive productivity and entertainment.",
	},
	"laptop": {
		"Thin and light chassis with all-day battery life and fast charging support.",
		"High-performance CPU and dedicated GPU for creative workloads and gaming.",
		"Stunning OLED display with 100% DCI-P3 color gamut for content creators.",
	},
	"headset": {
		"50mm drivers deliver deep bass and crisp highs for immersive audio.",
		"Noise-cancelling microphone with flip-to-mute and Discord-certified clarity.",
		"Memory foam ear cushions and adjustable headband for all-day comfort.",
	},
	"webcam": {
		"1080p 60fps streaming with autofocus and built-in stereo microphone.",
		"4K HDR sensor with background replacement and ring light compatibility.",
		"Wide-angle 90° FOV perfect for video calls, streaming, and recording.",
	},
}

func buildProducts() []model.CreateProductInput {
	var out []model.CreateProductInput

	// ~12-15 products per category = ~80 total
	for _, spec := range specs {
		count := 12 + rand.Intn(4) // 12-15
		for i := 0; i < count; i++ {
			brand := spec.brands[rand.Intn(len(spec.brands))]
			adj := adjectives[rand.Intn(len(adjectives))]
			model_ := 100 + rand.Intn(900)
			image := spec.images[rand.Intn(len(spec.images))]
			descs := descriptions[spec.name]
			desc := descs[rand.Intn(len(descs))]
			price := spec.minP + rand.Float64()*(spec.maxP-spec.minP)
			price = float64(int(price*100)) / 100 // round to 2dp
			stock := 5 + rand.Intn(96)            // 5-100

			out = append(out, model.CreateProductInput{
				Name:        fmt.Sprintf("%s %s %s %d", brand, adj, spec.name, model_),
				Brand:       brand,
				Description: desc,
				Price:       price,
				Stock:       stock,
				Category:    spec.name,
				Image:       image,
			})
		}
	}
	return out
}
