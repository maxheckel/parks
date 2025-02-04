package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/maxheckel/parks/database"
	"github.com/maxheckel/parks/models"
	"github.com/maxheckel/parks/services/maps"
	"os"
)

func main() {
	godotenv.Load(".env")
	os.Setenv("DB_HOST", "localhost")
	database.ConnectDb()
	mapsService, err := maps.New()
	if err != nil {
		panic(err)
	}
	parkModels := getParks()
	ctx := context.Background()
	for _, park := range parkModels {
		if park.Latitude > 0 {
			continue
		}
		placeID, lat, lng, err := mapsService.GetLocationLatAndLng(&ctx, park)
		if err != nil {
			panic(err)
		}
		park.Latitude = lat
		park.PlaceID = placeID
		park.Longitude = lng

		database.DB.Db.Save(&park)
		fmt.Println(park.ID)
	}
}

func getParks() []*models.Park {
	city := &models.City{
		Name:    "Columbus",
		Country: "United States",
		State:   "Ohio",
	}
	database.DB.Db.FirstOrCreate(&city)
	parkModels := []*models.Park{}
	for _, parkName := range parks {
		park := &models.Park{
			Name:   parkName,
			CityID: city.ID,
			City:   city,
		}
		database.DB.Db.Where("name = ?", park.Name).Where("city_id = ?", park.CityID).First(&park)
		if park.ID != 0 {
			parkModels = append(parkModels, park)
			continue
		}
		database.DB.Db.Create(&park)
		parkModels = append(parkModels, park)
	}
	return parkModels
}

var parks = []string{
	"Academy Park",
	"Albany Crossing Park",
	"Alexander AEP Park",
	"Alkire Woods Park",
	"Alum Crest Park",
	"Amvet Village Park",
	"Anheuser Busch Sports Park",
	"Antrim Park",
	"Argus Park",
	"Audubon Park",
	"Avalon Park",
	"Barnett Park",
	"Battelle Riverfront Park",
	"Beatty Park",
	"Beechcroft Park",
	"Beechwold Park",
	"Berliner Sports Park",
	"Berwick Park",
	"Bicentennial Park",
	"Big Run Park",
	"Big Walnut Park",
	"Blackburn Park",
	"Brandywine Park",
	"Brentnell Park",
	"Brevoort Park",
	"Brookside Woods Park",
	"Canini Park",
	"Carriage Place Park",
	"Cassady Park",
	"Casto Park",
	"Cedar Run Park",
	"Chaseland Park",
	"Cherry Bottom Park (South)",
	"City Gate Park",
	"Clinton-Como Park",
	"Clintonville Park",
	"Cody Park",
	"Columbus Commons (privately owned)",
	"Conner Park",
	"Cooke Park",
	"Cooper Park",
	"Crawford Farms Park",
	"Cremeans Park",
	"Deshler Park",
	"Devonshire Park",
	"Dexter Falls Park",
	"Dodge Park",
	"Dodge Skate Park",
	"Dorrian Commons Park (privately owned)",
	"Dorrian Green",
	"Driving Park",
	"Duranceau Park",
	"Easthaven Park",
	"Elk Run Park",
	"English Park",
	"Fairwood Park",
	"Flint Park",
	"Forest Park East Park",
	"Frank Fetch Memorial Park",
	"Franklin Park",
	"Franks Park",
	"Freedom Park",
	"Galloway Ridge Park",
	"Genoa Park",
	"Georgian Heights Park",
	"Glen Echo Park",
	"Glenview Park",
	"Glenwood Park",
	"Glick Park (O'shaughnessy Dam Overlook)",
	"Godown Road Park",
	"Goodale Park",
	"Granville Park",
	"Greene Countrie Park",
	"Griggs Reservoir Park",
	"Hamilton And Spring Portal Park",
	"Hanford Village Park",
	"Hard Road Park",
	"Harrison Park",
	"Harrison Smith Park",
	"Harrison West Park",
	"Hauntz Park",
	"Hayden Falls Park",
	"Hayden Park",
	"Haydens Crossing Park",
	"Heer Park",
	"Hellbranch Park",
	"Helsel Park",
	"Highbluffs Park",
	"Hilliard Green Park",
	"Hilltonia Park",
	"Holton Park",
	"Hoover Reservoir Park",
	"Huy Road Park",
	"Independence Village Park",
	"Indian Mound Park",
	"Indianola Park",
	"Innis Park",
	"Italian Village Park",
	"Iuka Park",
	"Jefferson Woods Park",
	"Jeffrey Scioto Park",
	"Joan Park",
	"Karns Park",
	"Keller Park",
	"Kelley Park",
	"Kenlawn Park",
	"Kenney Park",
	"Kingsrowe Park",
	"Kirkwood Park",
	"Kobacker Park",
	"Kraner Park",
	"Krumm Park",
	"Lazelle Woods Park",
	"Lehman Estates Park",
	"Lincoln Park",
	"Lindbergh Park",
	"Linden Park",
	"Linwood Park",
	"Liv Moor Park",
	"Livingston Park",
	"Lower Scioto Park",
	"Madison Mills Park",
	"Maloney Park",
	"Marie Moreland Park",
	"Marion Franklin Park",
	"Martin Luther King Park",
	"Martin Park",
	"Maybury Park",
	"Mayme Moore Park",
	"Maynard & Summit Park",
	"Mccoy Park",
	"McFerson Commons",
	"Mckinley Park",
	"Mifflin Park",
	"Millbrook Park",
	"Milo Grogan Park",
	"Mock Park",
	"Moeller Park",
	"Nafzger Park",
	"Nelson Park",
	"New Beginnings Park",
	"Noe Bixby Park",
	"North Bank Park",
	"North East Park",
	"Northcrest Park",
	"Northern Woods Park",
	"Northgate Park",
	"Northmoor Park",
	"Northtowne Park",
	"O' Shaughnessy Reservoir Park",
	"Ohio Police and Fire Memorial Park (privately owned)",
	"Olde Sawmill Park",
	"Overbrook Ravine Park",
	"Palsgrove Park",
	"Park of Roses (Whetstone Park)",
	"Parkridge Park",
	"Pingue Park",
	"Polaris Founder's Park",
	"Portman Park",
	"Prestwick Commons Park",
	"Pride Park",
	"Pump House Park",
	"Pumphrey Park",
	"Redick Park",
	"Remembrance Park (privately owned)",
	"Reynolds Crossing Park",
	"Rhodes Park",
	"Rickenbacker Park",
	"Riverbend Park",
	"Riverside Green Park",
	"Riverway Kiwanis Park",
	"Roosevelt Park",
	"Sader Park",
	"Sancus Park",
	"Saunders Park",
	"Sawyer Park",
	"Schiller Park",
	"Scioto Audubon Metro Park",
	"Scioto Greenlawn Dam Park",
	"Scioto Trail Park",
	"Scioto Woods Park",
	"Sensenbrenner Park",
	"Shady Lane Park",
	"Sharon Meadows Park",
	"Shrum Mound (owned by the Ohio History Connection)",
	"Shepard Park",
	"Side By Side Park",
	"Sills Park",
	"Smith Road Park",
	"South Side Settlement Heritage Park",
	"Southeast Lions Park",
	"Southgate Park",
	"Southwood Mileusnich Park",
	"Spindler Road Park",
	"Stephen Drive Park",
	"Stockbridge Park",
	"Stoneridge Park",
	"Strawberry Farms Park",
	"Summitview Park",
	"Sycamore Hills Park",
	"Tanager Woods Park",
	"The Promenade (Scioto Mile)",
	"Thompson Park",
	"Three Creeks Park",
	"Thurber Park",
	"Old Deaf School Park (Topiary Park)",
	"Trabue Woods Park",
	"Tuttle Park",
	"Walden Park",
	"Walnut Hill Park",
	"Walnut View Park",
	"Waltham Woods Park",
	"Washington Gladden Social Justice Park",
	"Webster Park",
	"Weinland Park",
	"Westchester Park",
	"Westgate Park",
	"Westmoor Park",
	"Wexford Green Park",
	"Wheeler Memorial Park",
	"Whetstone Park",
	"Williams Creek Park",
	"Willis Park",
	"Willow Creek Park",
	"Wilson Avenue Park",
	"Wilson Road Park",
	"Winchester Meadows Park",
	"Winding Creek Park",
	"Windsor Park",
	"Wolfe Park",
	"Woodbridge Green Park",
	"Woodward Park",
	"Worthington Hills Park",
	"Wrexham Park",
	"Wynstone Park",
	"Battelle Darby Creek",
	"Blacklick Woods",
	"Blacklick Woods Golf Course",
	"Blendon Woods",
	"Chestnut Ridge",
	"Clear Creek",
	"Glacier Ridge",
	"Greenway Trails",
	"Heritage Trail Park",
	"Highbanks",
	"Homestead",
	"Inniswood Metro Gardens",
	"Pickerington Ponds",
	"Prairie Oaks",
	"Quarry Trails",
	"Rocky Fork",
	"Scioto Audubon",
	"Scioto Grove",
	"Sharon Woods",
	"Slate Run",
	"Slate Run Farm",
	"Three Creeks",
	"Walnut Woods",
}
