import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/petifies/screens/petifies_create_screen/petifies_create_screen.dart';
import 'package:flutter/services.dart' show rootBundle;
import 'package:mobile/src/theme/themes.dart';

class LocationConfirmPage extends ConsumerStatefulWidget {
  const LocationConfirmPage({super.key});

  @override
  ConsumerState<LocationConfirmPage> createState() =>
      LocationConfirmPageState();
}

class LocationConfirmPageState extends ConsumerState<LocationConfirmPage> {
  late GoogleMapController _mapController;
  late String _darkMapStyle;
  late String _lightMapStyle;

  @override
  void initState() {
    super.initState();
    rootBundle.loadString(Constants.mapStylePath).then((string) {
      _lightMapStyle = string;
    });
    rootBundle.loadString(Constants.darkMapStylePath).then((string) {
      _darkMapStyle = string;
    });
  }

  @override
  Widget build(BuildContext context) {
    final address = ref.read(selectedAddressProvider);
    if (address == null) {
      return SizedBox.shrink();
    }
    final startLocation = LatLng(address.latitude, address.longitude);

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
                  "Is the pin in the right spot?",
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
              children: [
                GoogleMap(
                  zoomGesturesEnabled: true,
                  zoomControlsEnabled: false,
                  initialCameraPosition: CameraPosition(
                    target: startLocation,
                    zoom: 16.0,
                  ),
                  mapType: MapType.normal,
                  onMapCreated: (controller) {
                    setState(() {
                      _mapController = controller;
                      if (Theme.of(context).brightness == Brightness.light) {
                        _mapController.setMapStyle(_lightMapStyle);
                      } else {
                        _mapController.setMapStyle(_darkMapStyle);
                      }
                    });
                  },
                  onCameraMove: (CameraPosition cameraPosition) {
                    ref.read(selectedAddressProvider.notifier).updateCoordinate(
                        cameraPosition.target.latitude,
                        cameraPosition.target.longitude);
                  },
                ),
                Center(
                  //picker image on google map
                  child: Icon(
                    FontAwesomeIcons.locationDot,
                    size: 32,
                    color: Themes.blueColor,
                  ),
                ),
              ],
            ),
          )
        ],
      ),
    );
  }
}
