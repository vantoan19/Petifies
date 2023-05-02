// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:math';

import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart' show rootBundle;
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:geolocator/geolocator.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/petifies/controllers/petifies_controllers.dart';
import 'package:mobile/src/features/petifies/navigators/explore_navigator.dart';
import 'package:mobile/src/features/petifies/screens/petifies_details_screen.dart';
import 'package:mobile/src/models/petifies.dart';
import 'package:mobile/src/providers/petifies_providers.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/widgets/buttons/go_back_button.dart';
import 'package:mobile/src/widgets/buttons/go_back_root_button.dart';
import 'package:mobile/src/widgets/map/petifies_map.dart';
import 'package:mobile/src/widgets/petifies/petifies.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'petifies_explore_screen.g.dart';

final isMapViewProvider = StateProvider((ref) => false);
final selectedPetifiesType = StateProvider((ref) => PetifiesType.DOG_WALKING);

@riverpod
class PetifiesList extends _$PetifiesList {
  @override
  List<PetifiesModel> build(PetifiesType type, bool isMapConsumer) {
    return [];
  }

  void addPetifies(List<PetifiesModel> petifies) {
    state = [...state, ...petifies];
  }
}

class PetifiesExploreScreen extends ConsumerStatefulWidget {
  const PetifiesExploreScreen({super.key});

  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _PetifiesExploreScreenState();
  }
}

class _PetifiesExploreScreenState extends ConsumerState {
  double? _currentLongitude = null;
  double? _currentLatitude = null;
  double _zoom = 16;

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
        _currentLatitude = 0;
        _currentLongitude = 0;
        _zoom = 1;
      });
      return;
    }

    await Geolocator.getCurrentPosition(desiredAccuracy: LocationAccuracy.high)
        .then((Position position) {
      setState(() {
        _currentLatitude = position.latitude;
        _currentLongitude = position.longitude;
        _zoom = 16;
      });
    }).catchError((e) {
      debugPrint(e);
    });
  }

  @override
  void initState() {
    super.initState();

    _getCurrentPosition();
  }

  Widget _getTabBarView() {
    return TabBarView(
      physics: NeverScrollableScrollPhysics(),
      children: [
        _PetifiesList(
          type: PetifiesType.DOG_WALKING,
          currentLongitude: _currentLongitude!,
          currentLatitude: _currentLatitude!,
        ),
        _PetifiesList(
          type: PetifiesType.CAT_PLAYING,
          currentLongitude: _currentLongitude!,
          currentLatitude: _currentLatitude!,
        ),
        _PetifiesList(
          type: PetifiesType.DOG_SITTING,
          currentLongitude: _currentLongitude!,
          currentLatitude: _currentLatitude!,
        ),
        _PetifiesList(
          type: PetifiesType.CAT_SITTING,
          currentLongitude: _currentLongitude!,
          currentLatitude: _currentLatitude!,
        ),
        _PetifiesList(
          type: PetifiesType.DOG_ADOPTION,
          currentLongitude: _currentLongitude!,
          currentLatitude: _currentLatitude!,
        ),
        _PetifiesList(
          type: PetifiesType.CAT_ADOPTION,
          currentLongitude: _currentLongitude!,
          currentLatitude: _currentLatitude!,
        ),
      ],
    );
  }

  void _onTypeChange(int value) {
    PetifiesType type;
    switch (value) {
      case 0:
        type = PetifiesType.DOG_WALKING;
        break;
      case 1:
        type = PetifiesType.CAT_PLAYING;
        break;
      case 2:
        type = PetifiesType.DOG_SITTING;
        break;
      case 3:
        type = PetifiesType.CAT_SITTING;
        break;
      case 4:
        type = PetifiesType.DOG_ADOPTION;
        break;
      case 5:
        type = PetifiesType.CAT_ADOPTION;
        break;
      default:
        type = PetifiesType.DOG_WALKING;
    }
    ref.read(selectedPetifiesType.notifier).state = type;
  }

  @override
  Widget build(BuildContext context) {
    if (_currentLatitude == null || _currentLongitude == null) {
      return Center(child: CircularProgressIndicator());
    }

    final isMapView = ref.watch(isMapViewProvider);

    return DefaultTabController(
      initialIndex: 0,
      length: 6,
      child: Scaffold(
        appBar: AppBar(
          leading: const GoBackRootButton(),
          title: Text(
            "Petifies",
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w800,
            ),
          ),
          centerTitle: true,
          actions: [
            IconButton(
              iconSize: 26,
              onPressed: () {},
              icon: const Icon(Icons.settings),
            ),
            IconButton(
              iconSize: 26,
              padding: EdgeInsets.symmetric(horizontal: 30),
              onPressed: () {},
              icon: const Icon(Icons.chat_outlined),
            )
          ],
          toolbarHeight: 80,
          bottom: TabBar(
            isScrollable: true,
            indicatorColor: Themes.blueColor,
            indicatorWeight: 4,
            labelColor: Themes.blueColor,
            unselectedLabelColor: Theme.of(context).colorScheme.secondary,
            onTap: _onTypeChange,
            tabs: [
              Tab(
                icon: Icon(FontAwesomeIcons.dog),
                text: "Dog walking",
              ),
              Tab(
                icon: Icon(FontAwesomeIcons.cat),
                text: "Cat playing",
              ),
              Tab(
                icon: Icon(FontAwesomeIcons.shieldDog),
                text: "Dog sitting",
              ),
              Tab(
                icon: Icon(FontAwesomeIcons.shieldCat),
                text: "Cat sitting",
              ),
              Tab(
                icon: Icon(FontAwesomeIcons.bone),
                text: "Dog adoption",
              ),
              Tab(
                icon: Icon(FontAwesomeIcons.bowlFood),
                text: "Cat adoption",
              ),
            ],
          ),
        ),
        body: Stack(
          children: [
            AnimatedOpacity(
                opacity: isMapView ? 0 : 1,
                duration: const Duration(milliseconds: 100),
                child: IgnorePointer(
                  ignoring: isMapView ? true : false,
                  child: _getTabBarView(),
                )),
            AnimatedOpacity(
                opacity: isMapView ? 1 : 0,
                duration: const Duration(milliseconds: 100),
                child: IgnorePointer(
                  ignoring: isMapView ? false : true,
                  child: _PetifiesMapViewer(
                    currentLatitude: _currentLatitude!,
                    currentLongitude: _currentLongitude!,
                    zoom: _zoom,
                  ),
                )),
            Align(
              alignment: Alignment.bottomCenter,
              child: Padding(
                padding: const EdgeInsets.only(bottom: 24),
                child: const _MapSwitchingButton(),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _PetifiesList extends ConsumerWidget {
  final PetifiesType type;
  final double currentLongitude;
  final double currentLatitude;
  const _PetifiesList({
    required this.type,
    required this.currentLongitude,
    required this.currentLatitude,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final petifiesController = ref.watch(
      listPetifiesControllerProvider(
        parameters: ListPetifiesParameters(
          type: type,
          longitude: currentLongitude,
          latitude: currentLatitude,
          radius: 1e9,
          pageSize: 40,
          offset: 0,
          isMapConsumer: false,
        ),
      ),
    );
    final isLoading = petifiesController.isLoading;

    final petifies = ref.watch(petifiesListProvider(type, false));

    return Column(children: [
      Expanded(
        child: Stack(
          children: [
            if (petifies.length > 0)
              ListView.builder(
                itemBuilder: (context, index) => ProviderScope(
                  child: ProviderScope(
                    key: ObjectKey(petifies[index].id),
                    overrides: [
                      petifiesInfoProvider.overrideWithValue(petifies[index])
                    ],
                    child: const Petifies(),
                  ),
                ),
                itemCount: petifies.length,
              ),
            if (petifies.length == 0 && !isLoading)
              _EmptyItemPlaceholder(
                itemType: "petify",
              ),
          ],
        ),
      ),
      if (isLoading)
        Center(
          child: CircularProgressIndicator(),
        ),
    ]);
  }
}

class _PetifiesMapViewer extends ConsumerStatefulWidget {
  final double currentLongitude;
  final double currentLatitude;
  final double zoom;
  const _PetifiesMapViewer({
    required this.currentLongitude,
    required this.currentLatitude,
    required this.zoom,
  });

  @override
  ConsumerState<_PetifiesMapViewer> createState() => _PetifiesMapViewerState();
}

class _PetifiesMapViewerState extends ConsumerState<_PetifiesMapViewer> {
  late LatLng _curPosition;
  late double _curZoom;
  late GoogleMapController _mapController;
  late String _darkMapStyle;
  late String _lightMapStyle;
  late BitmapDescriptor _mapPin;

  @override
  void initState() {
    super.initState();

    _curPosition = LatLng(widget.currentLatitude, widget.currentLongitude);
    _curZoom = widget.zoom;
    rootBundle.loadString(Constants.mapStylePath).then((string) {
      _lightMapStyle = string;
    });
    rootBundle.loadString(Constants.darkMapStylePath).then((string) {
      _darkMapStyle = string;
    });
    BitmapDescriptor.fromAssetImage(
            ImageConfiguration(size: const Size(12, 12)), Constants.mapPin)
        .then((value) => _mapPin = value);
  }

  @override
  Widget build(BuildContext context) {
    final type = ref.watch(selectedPetifiesType);
    final radius = (500 * pow(2, 21 - _curZoom)).toDouble();
    final petifiesAsync = ref.watch(listPetifiesControllerProvider(
        parameters: ListPetifiesParameters(
            type: type,
            longitude: _curPosition.longitude,
            latitude: _curPosition.latitude,
            radius: radius,
            pageSize: 0,
            offset: 0,
            isMapConsumer: true)));

    final petifies = petifiesAsync.value;

    return GoogleMap(
      mapType: MapType.normal,
      initialCameraPosition: CameraPosition(
          target: LatLng(widget.currentLatitude, widget.currentLongitude),
          zoom: widget.zoom),
      myLocationEnabled: true,
      myLocationButtonEnabled: false,
      onMapCreated: (GoogleMapController controller) {
        setState(() {
          _mapController = controller;
          if (Theme.of(context).brightness == Brightness.light) {
            _mapController.setMapStyle(_lightMapStyle);
          } else {
            _mapController.setMapStyle(_darkMapStyle);
          }
        });
      },
      onCameraMove: (position) {
        _curPosition = position.target;
        _curZoom = position.zoom;
      },
      onCameraIdle: () {
        setState(() {
          _curPosition = _curPosition;
          _curZoom = _curZoom;
        });
      },
      markers: petifies != null
          ? Set.from(petifies.map((e) => Marker(
                markerId: MarkerId(e.id),
                icon: _mapPin,
                position: LatLng(e.address.latitude, e.address.longitude),
                onTap: () => Navigator.pushNamed(
                  context,
                  PetifiesExploreNavigatorRoutes.petifiesDetails,
                  arguments: PetifiesDetailsScreenArguments(
                    petifiesData: e,
                  ),
                ),
              )))
          : Set(),
    );
  }
}

class _MapSwitchingButton extends ConsumerWidget {
  const _MapSwitchingButton({
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final isMapView = ref.watch(isMapViewProvider);

    return ElevatedButton.icon(
      onPressed: () {
        ref.read(isMapViewProvider.notifier).state = !isMapView;
      },
      icon: Icon(
        isMapView ? FontAwesomeIcons.map : FontAwesomeIcons.listUl,
        color: Themes.whiteColor,
        size: 16,
      ),
      label: Text(
        isMapView ? "Map" : "Petifies",
        style: TextStyle(color: Themes.whiteColor),
      ),
      style: ButtonStyle(
          shape: MaterialStateProperty.all<RoundedRectangleBorder>(
            RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(50.0),
            ),
          ),
          padding: MaterialStatePropertyAll<EdgeInsetsGeometry>(
            EdgeInsets.symmetric(vertical: 10, horizontal: 18),
          ),
          backgroundColor: MaterialStatePropertyAll<Color>(
            Themes.blueColor90,
          )),
    );
  }
}

class _EmptyItemPlaceholder extends StatelessWidget {
  final String itemType;
  const _EmptyItemPlaceholder({
    Key? key,
    required this.itemType,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.only(top: 12),
      child: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          mainAxisSize: MainAxisSize.min,
          children: [
            SizedBox(
              width: 150,
              height: 150,
              child: Image.asset(Constants.emptyBoxPng),
            ),
            Text("No $itemType to show.")
          ],
        ),
      ),
    );
  }
}
