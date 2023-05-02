// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'petifies_controllers.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#custom-getters-and-methods');

/// @nodoc
mixin _$ListPetifiesParameters {
  petifiesModel.PetifiesType get type => throw _privateConstructorUsedError;
  double get longitude => throw _privateConstructorUsedError;
  double get latitude => throw _privateConstructorUsedError;
  double get radius => throw _privateConstructorUsedError;
  int get pageSize => throw _privateConstructorUsedError;
  int get offset => throw _privateConstructorUsedError;
  bool get isMapConsumer => throw _privateConstructorUsedError;

  @JsonKey(ignore: true)
  $ListPetifiesParametersCopyWith<ListPetifiesParameters> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $ListPetifiesParametersCopyWith<$Res> {
  factory $ListPetifiesParametersCopyWith(ListPetifiesParameters value,
          $Res Function(ListPetifiesParameters) then) =
      _$ListPetifiesParametersCopyWithImpl<$Res, ListPetifiesParameters>;
  @useResult
  $Res call(
      {petifiesModel.PetifiesType type,
      double longitude,
      double latitude,
      double radius,
      int pageSize,
      int offset,
      bool isMapConsumer});
}

/// @nodoc
class _$ListPetifiesParametersCopyWithImpl<$Res,
        $Val extends ListPetifiesParameters>
    implements $ListPetifiesParametersCopyWith<$Res> {
  _$ListPetifiesParametersCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? type = null,
    Object? longitude = null,
    Object? latitude = null,
    Object? radius = null,
    Object? pageSize = null,
    Object? offset = null,
    Object? isMapConsumer = null,
  }) {
    return _then(_value.copyWith(
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as petifiesModel.PetifiesType,
      longitude: null == longitude
          ? _value.longitude
          : longitude // ignore: cast_nullable_to_non_nullable
              as double,
      latitude: null == latitude
          ? _value.latitude
          : latitude // ignore: cast_nullable_to_non_nullable
              as double,
      radius: null == radius
          ? _value.radius
          : radius // ignore: cast_nullable_to_non_nullable
              as double,
      pageSize: null == pageSize
          ? _value.pageSize
          : pageSize // ignore: cast_nullable_to_non_nullable
              as int,
      offset: null == offset
          ? _value.offset
          : offset // ignore: cast_nullable_to_non_nullable
              as int,
      isMapConsumer: null == isMapConsumer
          ? _value.isMapConsumer
          : isMapConsumer // ignore: cast_nullable_to_non_nullable
              as bool,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$_ListPetifiesParametersCopyWith<$Res>
    implements $ListPetifiesParametersCopyWith<$Res> {
  factory _$$_ListPetifiesParametersCopyWith(_$_ListPetifiesParameters value,
          $Res Function(_$_ListPetifiesParameters) then) =
      __$$_ListPetifiesParametersCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {petifiesModel.PetifiesType type,
      double longitude,
      double latitude,
      double radius,
      int pageSize,
      int offset,
      bool isMapConsumer});
}

/// @nodoc
class __$$_ListPetifiesParametersCopyWithImpl<$Res>
    extends _$ListPetifiesParametersCopyWithImpl<$Res,
        _$_ListPetifiesParameters>
    implements _$$_ListPetifiesParametersCopyWith<$Res> {
  __$$_ListPetifiesParametersCopyWithImpl(_$_ListPetifiesParameters _value,
      $Res Function(_$_ListPetifiesParameters) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? type = null,
    Object? longitude = null,
    Object? latitude = null,
    Object? radius = null,
    Object? pageSize = null,
    Object? offset = null,
    Object? isMapConsumer = null,
  }) {
    return _then(_$_ListPetifiesParameters(
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as petifiesModel.PetifiesType,
      longitude: null == longitude
          ? _value.longitude
          : longitude // ignore: cast_nullable_to_non_nullable
              as double,
      latitude: null == latitude
          ? _value.latitude
          : latitude // ignore: cast_nullable_to_non_nullable
              as double,
      radius: null == radius
          ? _value.radius
          : radius // ignore: cast_nullable_to_non_nullable
              as double,
      pageSize: null == pageSize
          ? _value.pageSize
          : pageSize // ignore: cast_nullable_to_non_nullable
              as int,
      offset: null == offset
          ? _value.offset
          : offset // ignore: cast_nullable_to_non_nullable
              as int,
      isMapConsumer: null == isMapConsumer
          ? _value.isMapConsumer
          : isMapConsumer // ignore: cast_nullable_to_non_nullable
              as bool,
    ));
  }
}

/// @nodoc

class _$_ListPetifiesParameters implements _ListPetifiesParameters {
  _$_ListPetifiesParameters(
      {required this.type,
      required this.longitude,
      required this.latitude,
      required this.radius,
      required this.pageSize,
      required this.offset,
      required this.isMapConsumer});

  @override
  final petifiesModel.PetifiesType type;
  @override
  final double longitude;
  @override
  final double latitude;
  @override
  final double radius;
  @override
  final int pageSize;
  @override
  final int offset;
  @override
  final bool isMapConsumer;

  @override
  String toString() {
    return 'ListPetifiesParameters(type: $type, longitude: $longitude, latitude: $latitude, radius: $radius, pageSize: $pageSize, offset: $offset, isMapConsumer: $isMapConsumer)';
  }

  @override
  bool operator ==(dynamic other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$_ListPetifiesParameters &&
            (identical(other.type, type) || other.type == type) &&
            (identical(other.longitude, longitude) ||
                other.longitude == longitude) &&
            (identical(other.latitude, latitude) ||
                other.latitude == latitude) &&
            (identical(other.radius, radius) || other.radius == radius) &&
            (identical(other.pageSize, pageSize) ||
                other.pageSize == pageSize) &&
            (identical(other.offset, offset) || other.offset == offset) &&
            (identical(other.isMapConsumer, isMapConsumer) ||
                other.isMapConsumer == isMapConsumer));
  }

  @override
  int get hashCode => Object.hash(runtimeType, type, longitude, latitude,
      radius, pageSize, offset, isMapConsumer);

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$_ListPetifiesParametersCopyWith<_$_ListPetifiesParameters> get copyWith =>
      __$$_ListPetifiesParametersCopyWithImpl<_$_ListPetifiesParameters>(
          this, _$identity);
}

abstract class _ListPetifiesParameters implements ListPetifiesParameters {
  factory _ListPetifiesParameters(
      {required final petifiesModel.PetifiesType type,
      required final double longitude,
      required final double latitude,
      required final double radius,
      required final int pageSize,
      required final int offset,
      required final bool isMapConsumer}) = _$_ListPetifiesParameters;

  @override
  petifiesModel.PetifiesType get type;
  @override
  double get longitude;
  @override
  double get latitude;
  @override
  double get radius;
  @override
  int get pageSize;
  @override
  int get offset;
  @override
  bool get isMapConsumer;
  @override
  @JsonKey(ignore: true)
  _$$_ListPetifiesParametersCopyWith<_$_ListPetifiesParameters> get copyWith =>
      throw _privateConstructorUsedError;
}
