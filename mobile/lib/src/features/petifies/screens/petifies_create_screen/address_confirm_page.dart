import 'package:country_picker/country_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/petifies/controllers/google_map_controllers.dart';
import 'package:mobile/src/features/petifies/screens/petifies_create_screen/petifies_create_screen.dart';
import 'package:mobile/src/widgets/textfield/textformfield.dart';

class AddressConfirmPage extends ConsumerStatefulWidget {
  const AddressConfirmPage({super.key});

  @override
  ConsumerState<AddressConfirmPage> createState() => AddressConfirmPageState();
}

class AddressConfirmPageState extends ConsumerState<AddressConfirmPage> {
  final _formKey = GlobalKey<FormState>();
  final _streetController = TextEditingController();
  final _streetLineOneController = TextEditingController();
  final _cityController = TextEditingController();
  final _postalCodeController = TextEditingController();
  final _countryController = TextEditingController();

  @override
  void initState() {}

  @override
  Widget build(BuildContext context) {
    final selectedPlace = ref.read(selectedPlaceProvider);
    if (selectedPlace == null) {
      return SizedBox.shrink();
    }
    final address =
        ref.watch(getPlaceDetailControllerProvider(selectedPlace.placeId));

    return address.when(
      data: (data) {
        _streetController.text = data.street;
        _cityController.text = data.city;
        _postalCodeController.text = data.postalCode;
        _countryController.text = data.country;

        Future.delayed(const Duration(milliseconds: 500), () {
          ref.read(selectedAddressProvider.notifier).updateAddress(data);
          _formKey.currentState?.validate();
        });

        return SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Padding(
                padding: const EdgeInsets.fromLTRB(
                  Constants.petifiesExpoloreHorizontalPadding,
                  28,
                  Constants.petifiesExpoloreHorizontalPadding,
                  20,
                ),
                child: Text(
                  "Confirm your address",
                  style: TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 4),
                child: Form(
                  key: _formKey,
                  child: Column(
                    children: [
                      CustomTextFormField(
                        label: "Street",
                        icon: Icon(
                          FontAwesomeIcons.road,
                          size: 18,
                        ),
                        onChange: (value) {
                          ref
                              .read(selectedAddressProvider.notifier)
                              .updateStreet(value);
                          _formKey.currentState?.validate();
                        },
                        controller: _streetController,
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return 'Street is required';
                          }
                          return null;
                        },
                      ),
                      CustomTextFormField(
                        label: "Apartment, Building, etc. (Optional)",
                        icon: Icon(
                          FontAwesomeIcons.building,
                          size: 18,
                        ),
                        onChange: (value) {
                          ref
                              .read(selectedAddressProvider.notifier)
                              .updateAddressLineOne(value);
                        },
                        controller: _streetLineOneController,
                      ),
                      CustomTextFormField(
                        label: "City",
                        icon: Icon(
                          FontAwesomeIcons.city,
                          size: 18,
                        ),
                        onChange: (value) {
                          ref
                              .read(selectedAddressProvider.notifier)
                              .updateCity(value);
                          _formKey.currentState?.validate();
                        },
                        controller: _cityController,
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return 'City is required';
                          }
                          return null;
                        },
                      ),
                      CustomTextFormField(
                        label: "Postal code",
                        icon: Icon(
                          FontAwesomeIcons.envelope,
                          size: 18,
                        ),
                        onChange: (value) {
                          ref
                              .read(selectedAddressProvider.notifier)
                              .updatePostalCode(value);
                          _formKey.currentState?.validate();
                        },
                        controller: _postalCodeController,
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return 'Postal code is required';
                          }
                          return null;
                        },
                      ),
                      GestureDetector(
                        behavior: HitTestBehavior.translucent,
                        onTap: () => showCountryPicker(
                          context: context,
                          onSelect: (Country country) {
                            _countryController.text = country.name;
                          },
                        ),
                        child: IgnorePointer(
                          child: CustomTextFormField(
                            label: "Country",
                            icon: Icon(
                              FontAwesomeIcons.flag,
                              size: 18,
                            ),
                            suffixIcon: Icon(
                              FontAwesomeIcons.angleDown,
                              size: 18,
                            ),
                            onChange: (value) {
                              ref
                                  .read(selectedAddressProvider.notifier)
                                  .updateCountry(value);
                              _formKey.currentState?.validate();
                            },
                            controller: _countryController,
                            readOnly: true,
                            validator: (value) {
                              if (value == null || value.isEmpty) {
                                return 'Country is required';
                              }
                              return null;
                            },
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
        );
      },
      loading: () => CircularProgressIndicator(),
      error: (error, stackTrace) => Text('Error: $error'),
    );
  }
}
