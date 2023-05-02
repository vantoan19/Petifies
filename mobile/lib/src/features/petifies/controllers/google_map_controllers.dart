import 'dart:convert';

import 'package:http/http.dart';
import 'package:mobile/src/models/address.dart';

import 'package:mobile/src/models/suggestion.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:uuid/uuid.dart';

part 'google_map_controllers.g.dart';

@Riverpod(keepAlive: false)
class ListMapSuggestionsController extends _$ListMapSuggestionsController {
  static final String apikey = "AIzaSyA4b01Pv8MtYBiUYGJmhDDI8eDNtk9LwJo";

  Future<List<MapSuggestion>> _fetchSuggestions(String inp, String lang) async {
    final response = await Client().get(
        Uri.https("maps.googleapis.com", "/maps/api/place/autocomplete/json", {
      "input": inp,
      "types": "address",
      "language": lang,
      "key": apikey,
      "sessiontoken": this.sessionToken,
    }));

    if (response.statusCode == 200) {
      final result = json.decode(response.body);
      if (result['status'] == 'OK') {
        return result['predictions']
            .map<MapSuggestion>(
                (p) => MapSuggestion(p['place_id'], p['description']))
            .toList();
      }
      if (result['status'] == 'ZERO_RESULTS') {
        return [];
      }
      throw Exception(result['error_message']);
    } else {
      throw Exception('Failed to fetch suggestion');
    }
  }

  @override
  Future<List<MapSuggestion>> build(String sessionToken) async {
    return <MapSuggestion>[];
  }

  Future<List<MapSuggestion>> fetchSuggestions(String inp, String lang) async {
    return _fetchSuggestions(inp, lang);
  }
}

@Riverpod(keepAlive: false)
class GetPlaceDetailController extends _$GetPlaceDetailController {
  static final String apikey = "AIzaSyA4b01Pv8MtYBiUYGJmhDDI8eDNtk9LwJo";

  Future<Address> getDetails() async {
    final response = await Client()
        .get(Uri.https("maps.googleapis.com", "/maps/api/place/details/json", {
      "place_id": placeId,
      "fields": "address_component,geometry",
      "key": apikey,
      "sessiontoken": Uuid().v4(),
    }));

    if (response.statusCode == 200) {
      final result = json.decode(response.body);
      if (result['status'] == 'OK') {
        final components =
            result['result']['address_components'] as List<dynamic>;
        // build result
        final geometry = result['result']['geometry'];
        String? streetNumber;
        String? street;
        String? city;
        String? country;
        String? zipCode;

        components.forEach((c) {
          final List type = c['types'];
          if (type.contains('street_number')) {
            streetNumber = c['long_name'];
          }
          if (type.contains('route')) {
            street = c['long_name'];
          }
          if (type.contains('locality') ||
              type.contains("administrative_area_level_1")) {
            city = c['long_name'];
          }
          if (type.contains('country')) {
            country = c['long_name'];
          }
          if (type.contains('postal_code')) {
            zipCode = c['long_name'];
          }
        });
        return Address(
          street: (street != null ? street! : "") +
              (streetNumber != null ? ", $streetNumber" : ""),
          city: city ?? "",
          country: country ?? "",
          postalCode: zipCode ?? "",
          longitude: geometry['location']['lng'],
          latitude: geometry['location']['lat'],
        );
      }
      throw Exception(result['error_message']);
    } else {
      throw Exception('Failed to fetch place details');
    }
  }

  @override
  Future<Address> build(String placeId) async {
    return getDetails();
  }
}
