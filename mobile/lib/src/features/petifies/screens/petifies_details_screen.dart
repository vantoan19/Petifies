// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter/services.dart' show rootBundle;
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/petifies/controllers/petifies_session_controllers.dart';
import 'package:mobile/src/features/petifies/controllers/review_controllers.dart';
import 'package:mobile/src/features/petifies/screens/petifies_create_review_screen.dart';
import 'package:mobile/src/features/petifies/screens/petifies_create_session_screen.dart';
import 'package:mobile/src/models/petifies.dart';
import 'package:mobile/src/providers/petifies_providers.dart';
import 'package:mobile/src/providers/user_model_providers.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/widgets/appbars/petifies_details_appbar.dart';
import 'package:mobile/src/widgets/carousel_slider/carousel_slider_with_indicator.dart';
import 'package:mobile/src/widgets/petifies/petifies_session.dart';
import 'package:mobile/src/widgets/petifies/review_card.dart';
import 'package:mobile/src/widgets/status/status.dart';
import 'package:mobile/src/widgets/user_avatar/user_avatar.dart';

class PetifiesDetailsScreenArguments {
  final PetifiesModel petifiesData;
  PetifiesDetailsScreenArguments({
    required this.petifiesData,
  });
}

class PetifiesDetailsScreen extends ConsumerWidget {
  const PetifiesDetailsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final petifiesInfo = ref.watch(petifiesInfoProvider);
    final width = MediaQuery.of(context).size.width;

    return Scaffold(
      appBar: PetifiesDetailsAppBar(),
      body: SingleChildScrollView(
        child: Column(
          children: [
            ImageCarouselSliderWithIndicators(
              images: petifiesInfo.image,
              height: width * 0.6,
              width: width,
              roundCorner: false,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(
                  horizontal: Constants.petifiesExpoloreHorizontalPadding),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const _PetifiesHead(),
                  Divider(
                    thickness: 0.5,
                    height: 50,
                    color: Theme.of(context).colorScheme.secondary,
                  ),
                  const _PetAndOwnerInfo(),
                  Divider(
                    thickness: 0.5,
                    height: 50,
                    color: Theme.of(context).colorScheme.secondary,
                  ),
                  const _PetifiesDescription(),
                  Divider(
                    thickness: 0.5,
                    height: 50,
                    color: Theme.of(context).colorScheme.secondary,
                  ),
                  const _LocationDetail(),
                  Divider(
                    thickness: 0.5,
                    height: 50,
                    color: Theme.of(context).colorScheme.secondary,
                  ),
                  const _Sessions(),
                  Divider(
                    thickness: 0.5,
                    height: 50,
                    color: Theme.of(context).colorScheme.secondary,
                  ),
                  const _Reviews(),
                  Divider(
                    thickness: 0.5,
                    height: 50,
                    color: Theme.of(context).colorScheme.secondary,
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _PetifiesHead extends ConsumerWidget {
  const _PetifiesHead({super.key});

  Widget statusWidget(PetifiesStatus status) {
    switch (status) {
      case PetifiesStatus.AVAILABLE:
        return Status(
          color: Themes.blueColor,
          margin: EdgeInsets.symmetric(vertical: 8),
          label: "AVAILABLE",
          textSize: 12,
          textColor: Themes.whiteColor,
        );
      case PetifiesStatus.UNAVAILABLE:
        return Status(
          color: Themes.yellowColor,
          margin: EdgeInsets.symmetric(vertical: 8),
          label: "UNAVAILABLE",
          textSize: 12,
          textColor: Themes.whiteColor,
        );
      default:
        return Status(
          color: Themes.greyColor,
          margin: EdgeInsets.symmetric(vertical: 8),
          label: "UNKNOWN",
          textSize: 12,
          textColor: Themes.blackColor,
        );
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final title = ref.watch(petifiesInfoProvider.select((info) => info.title));
    final status =
        ref.watch(petifiesInfoProvider.select((info) => info.status));
    final address =
        ref.watch(petifiesInfoProvider.select((info) => info.address));

    return Padding(
      padding: const EdgeInsets.only(top: 20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            title,
            style: TextStyle(
              fontSize: 24,
              fontWeight: FontWeight.w400,
            ),
          ),
          statusWidget(status),
          Text(
            address.toString(),
            style: TextStyle(
              fontSize: 13,
              fontWeight: FontWeight.w300,
            ),
          ),
        ],
      ),
    );
  }
}

class _PetAndOwnerInfo extends ConsumerWidget {
  const _PetAndOwnerInfo({super.key});

  String invitation(PetifiesType type) {
    switch (type) {
      case PetifiesType.DOG_WALKING:
        return "Let's talk me for a walk!";
      case PetifiesType.CAT_PLAYING:
        return "Hey, let's play with me!";
      case PetifiesType.DOG_SITTING:
        return "I'm a good dog, please take care of me for a while";
      case PetifiesType.CAT_SITTING:
        return "Meow meow, would you like to take me for a while?";
      case PetifiesType.DOG_ADOPTION:
        return "Wolf wolf, would you like to be my owner?";
      case PetifiesType.CAT_ADOPTION:
        return "We're a good match, let me be your cat!";
      default:
        return "";
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final petName =
        ref.watch(petifiesInfoProvider.select((info) => info.petName));
    final owner = ref.watch(petifiesInfoProvider.select((info) => info.owner));
    final type = ref.watch(petifiesInfoProvider.select((info) => info.type));

    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Flexible(
          flex: 4,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                "${petName} owned by ${owner.firstName} ${owner.lastName}",
                style: TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.w400,
                ),
              ),
              Padding(
                padding: const EdgeInsets.only(top: 8.0),
                child: Text(
                  invitation(type),
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w300,
                    color: Theme.of(context).colorScheme.secondary,
                  ),
                ),
              ),
            ],
          ),
        ),
        Flexible(
          flex: 1,
          child: UserAvatar(
            userAvatar: owner.userAvatar,
          ),
        ),
      ],
    );
  }
}

class _PetifiesDescription extends ConsumerWidget {
  const _PetifiesDescription({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final description =
        ref.watch(petifiesInfoProvider.select((info) => info.description));

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          "About the Petifies",
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.w400,
          ),
        ),
        Padding(
          padding: const EdgeInsets.only(top: 12.0),
          child: Text(
            description,
            style: TextStyle(
              fontSize: 15,
              fontWeight: FontWeight.w300,
              color: Theme.of(context).colorScheme.secondary,
              height: 1.35,
            ),
          ),
        ),
      ],
    );
  }
}

class _LocationDetail extends ConsumerStatefulWidget {
  const _LocationDetail({super.key});
  @override
  ConsumerState<_LocationDetail> createState() => LocationConfirmPageState();
}

class LocationConfirmPageState extends ConsumerState<_LocationDetail> {
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
    final address =
        ref.watch(petifiesInfoProvider.select((value) => value.address));

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          "Where you can pick up my pet?",
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.w400,
          ),
        ),
        Padding(
          padding: const EdgeInsets.only(top: 24),
          child: ClipRRect(
            borderRadius: BorderRadius.circular(18),
            child: Container(
              height: 300,
              child: GoogleMap(
                zoomGesturesEnabled: true,
                initialCameraPosition: CameraPosition(
                    target: LatLng(address.latitude, address.longitude),
                    zoom: 14),
                zoomControlsEnabled: false,
                scrollGesturesEnabled: true,
                rotateGesturesEnabled: true,
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
                markers: {
                  Marker(
                    markerId: const MarkerId("currentLocation"),
                    position: LatLng(address.latitude, address.longitude),
                  )
                },
              ),
            ),
          ),
        ),
      ],
    );
    ;
  }
}

class _Sessions extends ConsumerWidget {
  const _Sessions({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final petifiesId =
        ref.watch(petifiesInfoProvider.select((value) => value.id));
    final petifiesOwner =
        ref.watch(petifiesInfoProvider.select((value) => value.owner));
    final sessions = ref.watch(
        listSessionsByPetifiesIdControllerProvider(petifiesId: petifiesId));

    final user = ref.watch(myUserProvider);
    final err = user.maybeWhen(
      error: (e, stackTrace) => e.toString(),
      orElse: () => "no err",
    );

    final userInfo = user.when(
      data: (data) => data,
      error: (e, stackTrace) => null,
      loading: () => null,
    );

    if (userInfo == null) {
      return Placeholder();
    }

    return sessions.when(
      data: (data) {
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              "Sessions",
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.w400,
              ),
            ),
            if (data.length > 0)
              Padding(
                padding: const EdgeInsets.only(top: 12),
                child: SizedBox(
                  height: 160,
                  child: ListView.builder(
                    scrollDirection: Axis.horizontal,
                    shrinkWrap: true,
                    itemBuilder: (context, index) {
                      return ProviderScope(
                        key: ObjectKey(data[index].id),
                        overrides: [
                          petifiesSessionProvider.overrideWithValue(data[index])
                        ],
                        child: PetifiesSession(
                          showListProposalButton:
                              (userInfo.id == petifiesOwner.id),
                        ),
                      );
                    },
                    itemCount: data.length,
                  ),
                ),
              ),
            if (data.length == 0)
              _EmptyItemPlaceholder(
                itemType: "session",
              ),
            if (userInfo.id == petifiesOwner.id)
              Padding(
                padding: const EdgeInsets.only(top: 8.0),
                child: OutlinedButton(
                  onPressed: () {
                    showModalBottomSheet(
                        context: context,
                        isScrollControlled: true,
                        useSafeArea: true,
                        barrierColor: Theme.of(context).scaffoldBackgroundColor,
                        builder: (context) {
                          return PetifiesCreateSessionScreen(
                              petifiesId: petifiesId);
                        });
                  },
                  style: OutlinedButton.styleFrom(
                      minimumSize: Size.fromHeight(44),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(8),
                      )),
                  child: Text("Create new Session for your Petifies"),
                ),
              ),
          ],
        );
      },
      error: (error, stackTrace) => Center(
        child: Text(error.toString()),
      ),
      loading: () => Center(child: CircularProgressIndicator()),
    );
  }
}

class _Reviews extends ConsumerWidget {
  const _Reviews({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final petifiesId =
        ref.watch(petifiesInfoProvider.select((value) => value.id));
    final reviews = ref.watch(
        listReviewsByPetifidesIdControllerProvider(petifiesId: petifiesId));

    return reviews.when(
      data: (data) {
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              "Reviews",
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.w400,
              ),
            ),
            if (data.length > 0)
              Padding(
                padding: const EdgeInsets.only(top: 12),
                child: SizedBox(
                  height: 160,
                  child: ListView.builder(
                    scrollDirection: Axis.horizontal,
                    shrinkWrap: true,
                    itemBuilder: (context, index) {
                      return ProviderScope(
                        key: ObjectKey(data[index].id),
                        overrides: [
                          reviewProvider.overrideWithValue(data[index]),
                        ],
                        child: const ReviewCard(),
                      );
                    },
                    itemCount: data.length,
                  ),
                ),
              ),
            if (data.length == 0)
              _EmptyItemPlaceholder(
                itemType: "review",
              ),
            Padding(
              padding: const EdgeInsets.only(top: 8.0),
              child: OutlinedButton(
                onPressed: () {
                  showModalBottomSheet(
                      context: context,
                      isScrollControlled: true,
                      useSafeArea: true,
                      barrierColor: Theme.of(context).scaffoldBackgroundColor,
                      builder: (context) {
                        return CreateReviewScreen(petifiesId: petifiesId);
                      });
                },
                style: OutlinedButton.styleFrom(
                    minimumSize: Size.fromHeight(44),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8),
                    )),
                child: Text("Write your review"),
              ),
            ),
          ],
        );
      },
      error: (error, stackTrace) => Center(
        child: Text(error.toString()),
      ),
      loading: () => Center(child: CircularProgressIndicator()),
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
        child: SizedBox(
          height: 160,
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            mainAxisSize: MainAxisSize.min,
            children: [
              SizedBox(
                width: 100,
                height: 100,
                child: Image.asset(Constants.emptyBoxPng),
              ),
              Text("No $itemType to show.")
            ],
          ),
        ),
      ),
    );
  }
}
