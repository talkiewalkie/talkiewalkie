//
//  Api.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 08/09/2021.
//

import CoreLocation
import Foundation

struct Api {
    static let url = "https://dev.talkiewalkie.app/"

    private static func get<T>(_ url: URL, completion: @escaping (T?, Error?) -> Void) where T: Codable {
        URLSession.shared.dataTask(with: url) { data, _, _ in
            guard let data = data else { return }
            DispatchQueue.main.async {
                do {
                    let output = try JSONDecoder().decode(T.self, from: data)
                    return completion(output, nil)
                } catch let err {
                    return completion(nil, err)
                }
            }
        }.resume()
    }

    // MARK: - WALKS

    struct WalkAuthor: Codable {
        let uuid: String
        let handle: String
    }

    struct WalksItem: Codable {
        let title: String
        let description: String
        let uuid: String
        let coverUrl: String
        let author: WalkAuthor
        let distanceFromPoint: Float
    }

    static func walks(offset: Int = 0, position: CLLocationCoordinate2D?, completion: @escaping ([Api.WalksItem]?, Error?) -> Void) {
        var url = URLComponents(string: "\(Api.url)/walks")!
        url.queryItems = [URLQueryItem(name: "offset", value: String(offset))]
        if let pos = position {
            url.queryItems!.append(contentsOf: [URLQueryItem(name: "lat", value: String(pos.latitude)), URLQueryItem(name: "lng", value: String(pos.longitude))])
        }
        Api.get(url.url!, completion: completion)
    }

    // MARK: - WALK

    struct Walk: Codable {
        let title: String
        let description: String
        let uuid: String
        let coverUrl: String
        let audioUrl: String
        let author: WalkAuthor
    }

    static func walk(_ uuid: String, completion: @escaping (Api.Walk?, Error?) -> Void) {
        let url = URLComponents(string: "\(Api.url)/walk/\(uuid)")!
        Api.get(url.url!, completion: completion)
    }
}
