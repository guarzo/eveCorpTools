package config

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
