package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

func init() {
	globalRegistry.RegisterMap("hz_hazard_pits", HazardPitsMap{})
}

type HazardPitsMap struct{}

func (m HazardPitsMap) ID() string {
	return "hz_hazard_pits"
}

func (m HazardPitsMap) Meta() Metadata {
	return Metadata{
		Name:        "hz_hazard_pits",
		Description: "A map that that fills in grid-like pattern of squares with pits filled with hazard sauce. Every N turns the pits will fill with another layer of sauce up to a maximum of 4 layers which last a few cycles, then the pits drain and the pattern repeats",
		Author:      "Battlesnake",
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  4,
		BoardSizes:  FixedSizes(Dimensions{11, 11}),
		Tags:        []string{TAG_FOOD_PLACEMENT, TAG_HAZARD_PLACEMENT, TAG_SNAKE_PLACEMENT},
	}
}

func (m HazardPitsMap) AddHazardPits(board *rules.BoardState, settings rules.Settings, editor Editor) {
	for x := 0; x < board.Width; x++ {
		for y := 0; y < board.Height; y++ {
			if x%2 == 1 && y%2 == 1 {
				point := rules.Point{X: x, Y: y}
				isStartPosition := false
				for _, startPos := range hazardPitStartPositions {
					if startPos == point {
						isStartPosition = true
					}
				}
				if !isStartPosition {
					editor.AddHazard(point)
				}
			}
		}
	}
}

func (m HazardPitsMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if !m.Meta().BoardSizes.IsAllowable(initialBoardState.Width, initialBoardState.Height) {
		return rules.RulesetError("This map can only be played on a 11x11 board")
	}

	if len(initialBoardState.Snakes) > len(hazardPitStartPositions) {
		return rules.ErrorTooManySnakes
	}

	rand := settings.GetRand(0)

	rand.Shuffle(len(hazardPitStartPositions), func(i int, j int) {
		hazardPitStartPositions[i], hazardPitStartPositions[j] = hazardPitStartPositions[j], hazardPitStartPositions[i]
	})
	snakeIDs := make([]string, 0, len(initialBoardState.Snakes))
	for _, snake := range initialBoardState.Snakes {
		snakeIDs = append(snakeIDs, snake.ID)
	}

	tempBoardState := rules.NewBoardState(initialBoardState.Width, initialBoardState.Height)
	tempBoardState.Snakes = make([]rules.Snake, len(snakeIDs))

	for i := 0; i < len(snakeIDs); i++ {
		tempBoardState.Snakes[i] = rules.Snake{
			ID:     snakeIDs[i],
			Health: rules.SnakeMaxHealth,
		}
	}

	for index, snake := range initialBoardState.Snakes {
		head := hazardPitStartPositions[index]
		err := rules.PlaceSnake(tempBoardState, snake.ID, []rules.Point{head, head, head})
		if err != nil {
			return err
		}
	}

	err := rules.PlaceFoodFixed(rand, tempBoardState)
	if err != nil {
		return err
	}

	// Copy food from temp board state
	for _, f := range tempBoardState.Food {
		editor.AddFood(f)
	}

	// Copy snakes from temp board state
	for _, snake := range tempBoardState.Snakes {
		editor.PlaceSnake(snake.ID, snake.Body, snake.Health)
	}

	return nil
}

func (m HazardPitsMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	err := StandardMap{}.UpdateBoard(lastBoardState, settings, editor)
	if err != nil {
		return err
	}

	// Cycle 0 - no hazards
	// Cycle 1 - 1 layer
	// Cycle 2 - 2 layers
	// Cycle 3 - 3 layers
	// Cycle 4-6 - 4 layers of hazards

	if lastBoardState.Turn%settings.RoyaleSettings.ShrinkEveryNTurns == 0 {
		// Is it time to update the hazards
		layers := (lastBoardState.Turn / settings.RoyaleSettings.ShrinkEveryNTurns) % 7
		if layers > 4 {
			layers = 4
		}

		editor.ClearHazards()

		// Add 1-4 layers of hazard pits depending on the cycle
		for n := 0; n < layers; n++ {
			m.AddHazardPits(lastBoardState, settings, editor)
		}
	}
	return nil
}

var hazardPitStartPositions = []rules.Point{
	{X: 1, Y: 1},
	{X: 9, Y: 1},
	{X: 1, Y: 9},
	{X: 9, Y: 9},
}