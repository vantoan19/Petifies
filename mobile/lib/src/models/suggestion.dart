class MapSuggestion {
  final String placeId;
  final String description;

  MapSuggestion(this.placeId, this.description);

  @override
  String toString() {
    return 'MapSuggestion(description: $description, placeId: $placeId)';
  }
}
