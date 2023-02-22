import 'dart:io';
import 'dart:convert';
import 'dart:typed_data';

import 'package:grpc/grpc.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:flutter/services.dart' show ByteData, rootBundle;

class GrpcFlutterClient {
  static ClientChannel? _clientChannel;

  static Future<ClientChannel> getClient() async {
    ByteData byteData = await rootBundle.load(Constants.caCertPath);
    Uint8List bytes = byteData.buffer.asUint8List();

    if (_clientChannel == null) {
      _clientChannel = ClientChannel(
        'localhost',
        port: 80,
        options: ChannelOptions(
          credentials: ChannelCredentials.secure(
            certificates: bytes,
            onBadCertificate: (certificate, host) => host == 'localhost:80',
          ),
          idleTimeout: Duration(minutes: 1),
        ),
      );
    }
    return _clientChannel!;
  }
}
