// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

class Tokens {
  final String sessionId;
  final String accessToken;
  final String refreshToken;
  final int accessTokenExpiresAt;
  final int refreshTokenExpiresAt;
  Tokens({
    required this.sessionId,
    required this.accessToken,
    required this.refreshToken,
    required this.accessTokenExpiresAt,
    required this.refreshTokenExpiresAt,
  });

  Tokens copyWith({
    String? sessionId,
    String? accessToken,
    String? refreshToken,
    int? accessTokenExpiresAt,
    int? refreshTokenExpiresAt,
  }) {
    return Tokens(
      sessionId: sessionId ?? this.sessionId,
      accessToken: accessToken ?? this.accessToken,
      refreshToken: refreshToken ?? this.refreshToken,
      accessTokenExpiresAt: accessTokenExpiresAt ?? this.accessTokenExpiresAt,
      refreshTokenExpiresAt:
          refreshTokenExpiresAt ?? this.refreshTokenExpiresAt,
    );
  }

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'sessionId': sessionId,
      'accessToken': accessToken,
      'refreshToken': refreshToken,
      'accessTokenExpiresAt': accessTokenExpiresAt,
      'refreshTokenExpiresAt': refreshTokenExpiresAt,
    };
  }

  factory Tokens.fromMap(Map<String, dynamic> map) {
    return Tokens(
      sessionId: map['sessionId'] as String,
      accessToken: map['accessToken'] as String,
      refreshToken: map['refreshToken'] as String,
      accessTokenExpiresAt: map['accessTokenExpiresAt'] as int,
      refreshTokenExpiresAt: map['refreshTokenExpiresAt'] as int,
    );
  }

  String toJson() => json.encode(toMap());

  factory Tokens.fromJson(String source) =>
      Tokens.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  String toString() {
    return 'Tokens(sessionId: $sessionId, accessToken: $accessToken, refreshToken: $refreshToken, accessTokenExpiresAt: $accessTokenExpiresAt, refreshTokenExpiresAt: $refreshTokenExpiresAt)';
  }

  @override
  bool operator ==(covariant Tokens other) {
    if (identical(this, other)) return true;

    return other.sessionId == sessionId &&
        other.accessToken == accessToken &&
        other.refreshToken == refreshToken &&
        other.accessTokenExpiresAt == accessTokenExpiresAt &&
        other.refreshTokenExpiresAt == refreshTokenExpiresAt;
  }

  @override
  int get hashCode {
    return sessionId.hashCode ^
        accessToken.hashCode ^
        refreshToken.hashCode ^
        accessTokenExpiresAt.hashCode ^
        refreshTokenExpiresAt.hashCode;
  }
}
