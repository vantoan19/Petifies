// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

class Address {
  final String addressLineOne;
  final String addressLineTwo;
  final String street;
  final String district;
  final String city;
  final String region;
  final String postalCode;
  final String country;
  final double longitude;
  final double latitude;
  Address({
    this.addressLineOne = "",
    this.addressLineTwo = "",
    this.street = "",
    this.district = "",
    this.city = "",
    this.region = "",
    this.postalCode = "",
    this.country = "",
    this.longitude = 0.0,
    this.latitude = 0.0,
  });

  Address copyWith({
    String? addressLineOne,
    String? addressLineTwo,
    String? street,
    String? district,
    String? city,
    String? region,
    String? postalCode,
    String? country,
    double? longitude,
    double? latitude,
  }) {
    return Address(
      addressLineOne: addressLineOne ?? this.addressLineOne,
      addressLineTwo: addressLineTwo ?? this.addressLineTwo,
      street: street ?? this.street,
      district: district ?? this.district,
      city: city ?? this.city,
      region: region ?? this.region,
      postalCode: postalCode ?? this.postalCode,
      country: country ?? this.country,
      longitude: longitude ?? this.longitude,
      latitude: latitude ?? this.latitude,
    );
  }

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'addressLineOne': addressLineOne,
      'addressLineTwo': addressLineTwo,
      'street': street,
      'district': district,
      'city': city,
      'region': region,
      'postalCode': postalCode,
      'country': country,
      'longitude': longitude,
      'latitude': latitude,
    };
  }

  factory Address.fromMap(Map<String, dynamic> map) {
    return Address(
      addressLineOne: map['addressLineOne'] as String,
      addressLineTwo: map['addressLineTwo'] as String,
      street: map['street'] as String,
      district: map['district'] as String,
      city: map['city'] as String,
      region: map['region'] as String,
      postalCode: map['postalCode'] as String,
      country: map['country'] as String,
      longitude: map['longitude'] as double,
      latitude: map['latitude'] as double,
    );
  }

  String toJson() => json.encode(toMap());

  factory Address.fromJson(String source) =>
      Address.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  String toString() {
    return '${street}, ${city}, ${country}, ${postalCode}';
  }

  @override
  bool operator ==(covariant Address other) {
    if (identical(this, other)) return true;

    return other.addressLineOne == addressLineOne &&
        other.addressLineTwo == addressLineTwo &&
        other.street == street &&
        other.district == district &&
        other.city == city &&
        other.region == region &&
        other.postalCode == postalCode &&
        other.country == country &&
        other.longitude == longitude &&
        other.latitude == latitude;
  }

  @override
  int get hashCode {
    return addressLineOne.hashCode ^
        addressLineTwo.hashCode ^
        street.hashCode ^
        district.hashCode ^
        city.hashCode ^
        region.hashCode ^
        postalCode.hashCode ^
        country.hashCode ^
        longitude.hashCode ^
        latitude.hashCode;
  }
}
