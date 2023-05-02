// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/features/media/controllers/image_controller.dart';
import 'package:mobile/src/features/petifies/controllers/petifies_controllers.dart';
import 'package:mobile/src/features/petifies/screens/petifies_create_screen/address_confirm_page.dart';
import 'package:mobile/src/features/petifies/screens/petifies_create_screen/image_picker_page.dart';
import 'package:mobile/src/features/petifies/screens/petifies_create_screen/location_confirm_page.dart';
import 'package:mobile/src/features/petifies/screens/petifies_create_screen/map_selection_page.dart';
import 'package:mobile/src/features/petifies/screens/petifies_create_screen/petifies_info_form_page.dart';
import 'package:mobile/src/features/petifies/screens/petifies_create_screen/petifies_type_selection_page.dart';
import 'package:mobile/src/models/address.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/petifies.dart';
import 'package:mobile/src/models/suggestion.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/utils/navigation.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'petifies_create_screen.g.dart';

final petifiesTypeProvider = StateProvider<PetifiesType?>((ref) => null);
final selectedPlaceProvider = StateProvider<MapSuggestion?>((ref) => null);
final titleProvider = StateProvider<String?>((ref) => null);
final descriptionProvider = StateProvider<String?>((ref) => null);
final petNameProvider = StateProvider<String?>((ref) => null);

@Riverpod(keepAlive: true)
class SelectedAddress extends _$SelectedAddress {
  @override
  Address? build() {
    return null;
  }

  void updateAddress(Address address) {
    state = address;
  }

  void updateStreet(String street) {
    state = state?.copyWith(street: street);
  }

  void updateAddressLineOne(String addressLineOne) {
    state = state?.copyWith(addressLineOne: addressLineOne);
  }

  void updateCity(String city) {
    state = state?.copyWith(city: city);
  }

  void updateCountry(String country) {
    state = state?.copyWith(country: country);
  }

  void updatePostalCode(String postalCode) {
    state = state?.copyWith(postalCode: postalCode);
  }

  void updateCoordinate(double latitude, double longitude) {
    state = state?.copyWith(latitude: latitude, longitude: longitude);
  }
}

@Riverpod(keepAlive: true)
class ImageFiles extends _$ImageFiles {
  @override
  List<File> build() {
    return [];
  }

  void addImages(List<File> images) {
    state = [...state, ...images];
  }

  void removeImage(int index) {
    if (index >= state.length || index < 0) {
      return;
    }

    state = [...state.sublist(0, index), ...state.sublist(index + 1)];
  }
}

@Riverpod(keepAlive: true)
class ImageFutures extends _$ImageFutures {
  @override
  List<Future<Either<Failure, NetworkImageModel>>> build() {
    return [];
  }

  void addImages(List<Future<Either<Failure, NetworkImageModel>>> images) {
    state = [...state, ...images];
  }

  void removeImage(int index) async {
    if (index >= state.length || index < 0) {
      return;
    }

    final future = state[index];
    (await future).fold((l) => null, (r) {
      ref.read(imageControllerProvider.notifier).removeImage(uri: r.uri);
    });

    state = [...state.sublist(0, index), ...state.sublist(index + 1)];
  }
}

class CreatePetifiesScreen extends ConsumerStatefulWidget {
  const CreatePetifiesScreen({super.key});

  @override
  ConsumerState<CreatePetifiesScreen> createState() =>
      _CreatePetifiesScreenState();
}

class _CreatePetifiesScreenState extends ConsumerState<CreatePetifiesScreen> {
  PageController _pageController = PageController();
  int _currentPage = 0;

  void _submitCreatePetifies() {
    final type = ref.read(petifiesTypeProvider);
    final title = ref.read(titleProvider);
    final description = ref.read(descriptionProvider);
    final petName = ref.read(petNameProvider);
    final imgFutures = ref.read(imageFuturesProvider);
    final address = ref.read(selectedAddressProvider);

    ref.read(createPetifiesControllerProvider.notifier).createPetifies(
        type: type!,
        title: title!,
        description: description!,
        petName: petName!,
        imageFutures: imgFutures,
        address: address!);

    NavigatorUtil.goBack(context);
  }

  void _nextPage() {
    if (_currentPage == 5) {
      _submitCreatePetifies();
      return;
    }
    if (!_isNextDisable()) {
      setState(() {
        _currentPage++;
        _pageController.animateToPage(
          _currentPage,
          duration: Duration(
            milliseconds: 300,
          ),
          curve: Curves.linear,
        );
      });
    }
  }

  void _previousPage() {
    if (_currentPage > 0) {
      setState(() {
        _currentPage--;
        _pageController.animateToPage(
          _currentPage,
          duration: Duration(
            milliseconds: 300,
          ),
          curve: Curves.linear,
        );
      });
    }
  }

  bool _isDisableBack() {
    return _currentPage == 0;
  }

  bool _isNextDisable() {
    switch (_currentPage) {
      case 0:
        return ref.watch(petifiesTypeProvider) == null;
      case 1:
        return ref.watch(selectedPlaceProvider) == null;
      case 2:
        final address = ref.watch(selectedAddressProvider);
        if ((address?.street == null || address?.street == "") ||
            (address?.city == null || address?.city == "") ||
            (address?.country == null || address?.country == "") ||
            (address?.postalCode == null || address?.postalCode == "") ||
            (address?.latitude == null) ||
            (address?.longitude == null)) {
          return true;
        }
        break;
      case 4:
        final title = ref.watch(titleProvider);
        final description = ref.watch(descriptionProvider);
        if (PetifiesInfoFormPage.titleValidator(title) != null ||
            PetifiesInfoFormPage.descriptionValidator(description) != null) {
          return true;
        }
        break;
    }
    return false;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          onPressed: () {},
          icon: const Icon(Icons.close),
        ),
      ),
      body: Column(
        children: [
          Expanded(
            child: PageView(
              controller: _pageController,
              physics: NeverScrollableScrollPhysics(),
              children: [
                const PetifiesTypeSelectionPage(),
                const MapSelectionPage(),
                const AddressConfirmPage(),
                const LocationConfirmPage(),
                const PetifiesInfoFormPage(),
                const ImagePickerPage(),
              ],
            ),
          ),
          Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              LinearProgressIndicator(
                value: _currentPage / 5,
                minHeight: 6,
                color: Theme.of(context).colorScheme.inversePrimary,
                backgroundColor: Theme.of(context).colorScheme.tertiary,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(
                  horizontal: Constants.petifiesExpoloreHorizontalPadding,
                  vertical: 16,
                ),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    TextButton(
                      onPressed: _isDisableBack() ? null : _previousPage,
                      child: Text(
                        "Back",
                        style: TextStyle(
                          color: Theme.of(context).colorScheme.secondary,
                          fontSize: 16,
                          letterSpacing: 0.25,
                          fontWeight: FontWeight.w500,
                          decoration: TextDecoration.underline,
                          decorationThickness: 1,
                        ),
                      ),
                    ),
                    ElevatedButton(
                      onPressed: _isNextDisable() ? null : _nextPage,
                      child: Text(
                        _currentPage < 5 ? "Next" : "Create",
                        style: TextStyle(
                          color: Themes.whiteColor,
                          fontSize: 16,
                          letterSpacing: 0.25,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                      style: ElevatedButton.styleFrom(
                        padding: EdgeInsets.symmetric(
                          vertical: 16,
                          horizontal: 28,
                        ),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                      ),
                    ),
                  ],
                ),
              )
            ],
          ),
        ],
      ),
    );
  }
}
