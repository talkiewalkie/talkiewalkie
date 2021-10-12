//
//  Photos+authorization.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import Photos

extension PHPhotoLibrary {
    enum AuthorizationStatus {
        case justDenied
        case alreadyDenied
        case restricted
        case justLimited
        case alreadyLimited
        case justAuthorized
        case alreadyAuthorized
        case unknown

        var superStatus: PHAuthorizationStatus {
            switch self {
            case .justDenied, .alreadyDenied:
                return .denied
            case .restricted:
                return .restricted
            case .justLimited, .alreadyLimited:
                return .limited
            case .justAuthorized, .alreadyAuthorized:
                return .authorized
            default:
                return .notDetermined
            }
        }
    }

    class func checkAuthorization(completion: ((AuthorizationStatus?) -> Void)?) {
        let status = PHPhotoLibrary.authorizationStatus()
        switch status {
        case .authorized:
            completion?(.alreadyAuthorized)
        case .limited:
            completion?(.alreadyLimited)
        case .denied:
            completion?(.alreadyDenied)
        case .restricted:
            completion?(.restricted)
        default:
            completion?(.unknown)
        }
    }

    class func authorize(completion: ((AuthorizationStatus) -> Void)?) {
        let status = PHPhotoLibrary.authorizationStatus()
        switch status {
        case .authorized:
            completion?(.alreadyAuthorized)
        case .limited:
            completion?(.alreadyLimited)
        case .denied:
            completion?(.alreadyDenied)
        case .restricted:
            completion?(.restricted)
        case .notDetermined:
            PHPhotoLibrary.requestAuthorization { status in

                DispatchQueue.main.async {
                    switch status {
                    case .authorized:
                        completion?(.justAuthorized)
                    case .limited:
                        completion?(.justLimited)
                    case .denied:
                        completion?(.justDenied)
                    case .restricted:
                        completion?(.restricted)
                    case .notDetermined:
                        completion?(.unknown)
                    @unknown default:
                        completion?(.unknown)
                    }
                }
            }
        @unknown default:
            completion?(.unknown)
        }
    }
}
