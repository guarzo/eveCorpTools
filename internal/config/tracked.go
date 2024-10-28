package config

// CorporationIDs are the IDs of the corporations
var CorporationIDs = []int{98648442, 98730557, 98763685, 98743419}

// AllianceIDs are the IDs of the alliances
var AllianceIDs = []int{99010452}

// CharacterIDs are the IDs of the characters
var CharacterIDs = []int{1959376155, 2121524689, 96180548, 2118868995, 2118016167, 2114311509, 537223062, 2115754172, 629507683, 640170087, 2119887294, 1406208348, 1872552403, 2112148425, 404850015}

var ExcludeCharacters = []int{2116875456, 2120850653, 2121334187, 2121355778}

func TrackedCharacterID(characterID int) bool {
	for _, id := range CharacterIDs {
		if id == characterID {
			return true
		}
	}
	return false
}

func TrackedAllianceID(allianceID int) bool {
	for _, id := range AllianceIDs {
		if id == allianceID {
			return true
		}
	}
	return false
}

func TrackedCorporationID(corporationID int) bool {
	for _, id := range CorporationIDs {
		if id == corporationID {
			return true
		}
	}
	return false
}

func ExcludeCharacterID(characterID int) bool {
	for _, id := range ExcludeCharacters {
		if id == characterID {
			return true
		}
	}
	return false
}

func DisplayCharacter(characterID, corporationID, allianceID int) bool {
	return !ExcludeCharacterID(characterID) &&
		(TrackedCharacterID(characterID) || TrackedCorporationID(corporationID) || TrackedAllianceID(allianceID))
}
