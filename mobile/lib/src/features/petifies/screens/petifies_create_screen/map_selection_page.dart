import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:geolocator/geolocator.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/petifies/controllers/google_map_controllers.dart';
import 'package:mobile/src/features/petifies/screens/petifies_create_screen/petifies_create_screen.dart';
import 'package:mobile/src/models/address.dart';
import 'package:flutter/services.dart' show rootBundle;
import 'package:mobile/src/models/suggestion.dart';
import 'package:uuid/uuid.dart';

class MapSelectionPage extends ConsumerStatefulWidget {
  const MapSelectionPage();

  @override
  ConsumerState<MapSelectionPage> createState() => MapSelectionPageState();
}

class MapSelectionPageState extends ConsumerState<MapSelectionPage> {
  late GoogleMapController _mapController;
  LatLng? _geometryLocation;
  double? _zoom;
  late String _darkMapStyle;
  late String _lightMapStyle;
  String location = "Enter your location";

  Future<bool> _handleLocationPermission() async {
    bool serviceEnabled;
    LocationPermission permission;

    serviceEnabled = await Geolocator.isLocationServiceEnabled();
    if (!serviceEnabled) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(
          content: Text(
              'Location services are disabled. Please enable the services')));
      return false;
    }
    permission = await Geolocator.checkPermission();
    if (permission == LocationPermission.denied) {
      permission = await Geolocator.requestPermission();
      if (permission == LocationPermission.denied) {
        ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Location permissions are denied')));
        return false;
      }
    }
    if (permission == LocationPermission.deniedForever) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(
          content: Text(
              'Location permissions are permanently denied, we cannot request permissions.')));
      return false;
    }
    return true;
  }

  Future<void> _getCurrentPosition() async {
    final hasPermission = await _handleLocationPermission();
    if (!hasPermission) {
      setState(() {
        _geometryLocation = LatLng(0, 0);
        _zoom = 1;
      });
      return;
    }

    await Geolocator.getCurrentPosition(desiredAccuracy: LocationAccuracy.high)
        .then((Position position) {
      setState(() {
        _geometryLocation = LatLng(position.latitude, position.longitude);
        _zoom = 16;
      });
    }).catchError((e) {
      debugPrint(e);
    });
  }

  @override
  void initState() {
    super.initState();
    rootBundle.loadString(Constants.mapStylePath).then((string) {
      _lightMapStyle = string;
    });
    rootBundle.loadString(Constants.darkMapStylePath).then((string) {
      _darkMapStyle = string;
    });

    _getCurrentPosition();
  }

  @override
  Widget build(BuildContext context) {
    if (_geometryLocation == null) {
      return Center(child: CircularProgressIndicator());
    }

    final selectedPlace = ref.watch(selectedPlaceProvider);
    final String desc = selectedPlace?.description ?? "Enter your location";
    final address = (selectedPlace != null)
        ? ref.watch(getPlaceDetailControllerProvider(selectedPlace.placeId))
        : AsyncValue.data(Address(
            longitude: _geometryLocation!.longitude,
            latitude: _geometryLocation!.latitude));

    return Scaffold(
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(
              Constants.petifiesExpoloreHorizontalPadding,
              28,
              Constants.petifiesExpoloreHorizontalPadding,
              20,
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  "Where they can pick up your pet?",
                  style: TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.only(top: 8.0),
                  child: Text(
                    "Your full address is only shared with proposers after you accept their proposals.",
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w300,
                      color: Theme.of(context).colorScheme.secondary,
                    ),
                  ),
                ),
              ],
            ),
          ),
          Expanded(
            child: Stack(
              fit: StackFit.expand,
              children: [
                address.when(
                  data: (data) {
                    return GoogleMap(
                      zoomGesturesEnabled: false,
                      initialCameraPosition: CameraPosition(
                          target: LatLng(data.latitude, data.longitude),
                          zoom: 14),
                      zoomControlsEnabled: false,
                      scrollGesturesEnabled: false,
                      rotateGesturesEnabled: false,
                      mapType: MapType.normal,
                      onMapCreated: (controller) {
                        setState(() {
                          _mapController = controller;
                          if (Theme.of(context).brightness ==
                              Brightness.light) {
                            _mapController.setMapStyle(_lightMapStyle);
                          } else {
                            _mapController.setMapStyle(_darkMapStyle);
                          }
                        });
                      },
                      markers: {
                        Marker(
                          markerId: const MarkerId("currentLocation"),
                          position: LatLng(data.latitude, data.longitude),
                        )
                      },
                    );
                  },
                  error: (error, stackTrace) => Text("$error"),
                  loading: () => Center(child: CircularProgressIndicator()),
                ),
                Align(
                  alignment: Alignment.topCenter,
                  child: Padding(
                    padding: const EdgeInsets.only(top: 24),
                    child: InkWell(
                      onTap: () async {
                        final sessionToken = Uuid().v4();
                        final MapSuggestion? result = await showSearch(
                          context: context,
                          delegate: _AddressSearch(
                            sessionToken: sessionToken,
                            ref: ref,
                          ),
                        );

                        if (result != null && result.placeId != "") {
                          ref.read(selectedPlaceProvider.notifier).state =
                              result;
                        }
                      },
                      child: Card(
                        color: Theme.of(context).colorScheme.tertiary,
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.all(Radius.circular(50)),
                          side: BorderSide(
                            width: 1,
                            color: Theme.of(context).colorScheme.inversePrimary,
                          ),
                        ),
                        child: Container(
                            padding: EdgeInsets.all(4),
                            width: MediaQuery.of(context).size.width - 60,
                            child: ListTile(
                              title: Text(
                                desc,
                                style: TextStyle(
                                    fontSize: 16, fontWeight: FontWeight.w600),
                              ),
                              leading: Icon(
                                Icons.location_on,
                                size: 18,
                              ),
                              dense: true,
                            )),
                      ),
                    ),
                  ),
                )
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class SelectedPlaceDescription extends ConsumerWidget {
  const SelectedPlaceDescription({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final selectedPlace = ref.watch(selectedPlaceProvider);
    final String desc = selectedPlace?.description ?? "Enter your location";

    return Container(
        padding: EdgeInsets.all(4),
        width: MediaQuery.of(context).size.width - 60,
        child: ListTile(
          title: Text(
            desc,
            style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
          ),
          leading: Icon(
            Icons.location_on,
            size: 18,
          ),
          dense: true,
        ));
  }
}

class _AddressSearch extends SearchDelegate<MapSuggestion> {
  final String sessionToken;
  final WidgetRef ref;
  _AddressSearch({
    required this.sessionToken,
    required this.ref,
  });

  @override
  ThemeData appBarTheme(BuildContext context) {
    return Theme.of(context).copyWith(
      appBarTheme: Theme.of(context).appBarTheme.copyWith(),
    );
  }

  @override
  List<Widget> buildActions(BuildContext context) {
    return [
      IconButton(
        tooltip: 'Clear',
        icon: Icon(Icons.clear),
        onPressed: () {
          query = '';
        },
      )
    ];
  }

  @override
  Widget buildLeading(BuildContext context) {
    return IconButton(
      tooltip: 'Back',
      icon: Icon(Icons.arrow_back),
      onPressed: () {
        close(context, MapSuggestion("", ""));
      },
    );
  }

  @override
  Widget buildResults(BuildContext context) {
    return SizedBox.shrink();
  }

  @override
  Widget buildSuggestions(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(
          horizontal: Constants.petifiesExpoloreHorizontalPadding),
      child: Center(
        child: FutureBuilder(
          future: query == ""
              ? null
              : ref
                  .read(listMapSuggestionsControllerProvider(sessionToken)
                      .notifier)
                  .fetchSuggestions(
                      query, Localizations.localeOf(context).languageCode),
          builder: (context, snapshot) => query == ''
              ? SizedBox.shrink()
              : snapshot.hasData
                  ? ListView.builder(
                      itemBuilder: (context, index) => ListTile(
                        contentPadding: EdgeInsets.zero,
                        leading: Icon(
                          Icons.location_city_outlined,
                          size: 20,
                        ),
                        title: Text((snapshot.data?[index] as MapSuggestion)
                            .description),
                        onTap: () {
                          close(
                              context, snapshot.data?[index] as MapSuggestion);
                        },
                      ),
                      itemCount: snapshot.data?.length,
                    )
                  : CircularProgressIndicator(),
        ),
      ),
    );
  }
}
