import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/models/petifies.dart';
import 'package:mobile/src/models/petifies_proposal.dart';
import 'package:mobile/src/models/petifies_session.dart';
import 'package:mobile/src/models/review.dart';

final petifiesInfoProvider =
    Provider<PetifiesModel>((ref) => throw UnimplementedError());

final petifiesSessionProvider =
    Provider<PetifiesSessionModel>((ref) => throw UnimplementedError());

final reviewProvider =
    Provider<ReviewModel>((ref) => throw UnimplementedError());

final proposalProvider =
    Provider<PetifiesProposalModel>((ref) => throw UnimplementedError());
