//
//  Api.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 08/09/2021.
//

import CoreLocation
import Foundation

struct Api {
    #if DEBUG
    static let url = "https://dev.talkiewalkie.app"
    #else
    static let url = "https://api.talkiewalkie.app"
    #endif
    
    var token: String
    
    private func get<T>(_ url: URL, completion: @escaping (T?, Error?) -> Void) where T: Codable {
        var request = URLRequest(url: url)
        request.addValue("Bearer \(token)", forHTTPHeaderField: "X-TalkieWalkie-Auth")
        
        print("\(Date()) GET '\(url)'")
        URLSession.shared.dataTask(with: request) { data, res, httpErr in
            if let err = httpErr { print(err); return }
            if let httpResponse = res as? HTTPURLResponse {
                if httpResponse.statusCode > 299 {
                    let body = data != nil ? String(data: data!, encoding: .utf8)! : "empty body"
                    print("\(Date()) GET '\(url)' -> \(httpResponse.statusCode): '\(body)'")
                    
                    return
                }
            }
            guard let data = data else {
                print("\(Date()) GET '\(url)' -> no data")
                
                return
            }
            
            DispatchQueue.main.async {
                do {
                    let output = try JSONDecoder().decode(T.self, from: data)
                    
                    return completion(output, nil)
                } catch let err {
                    print("\(Date()) GET '\(url)' Could not decode JSON: \(err)")
                    
                    return completion(nil, err)
                }
            }
        }.resume()
    }
    
    private func post<I, T>(_ url: URL, payload: I, completion: @escaping (T?, Error?) -> Void) where T: Codable, I: Codable {
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.addValue("Bearer \(token)", forHTTPHeaderField: "X-TalkieWalkie-Auth")
        let jsonData = try? JSONEncoder().encode(payload)
        request.httpBody = jsonData
        
        print("\(Date()) GET '\(url)'")
        URLSession.shared.dataTask(with: request) { data, res, httpErr in
            if let err = httpErr { print(err); return }
            if let httpResponse = res as? HTTPURLResponse {
                if httpResponse.statusCode > 299 {
                    let body = data != nil ? String(data: data!, encoding: .utf8)! : "empty body"
                    print("\(Date()) GET '\(url)' -> \(httpResponse.statusCode): '\(body)'")
                    
                    return
                }
            }
            guard let data = data else {
                print("\(Date()) GET '\(url)' -> no data")
                
                return
            }
            
            DispatchQueue.main.async {
                do {
                    let output = try JSONDecoder().decode(T.self, from: data)
                    
                    return completion(output, nil)
                } catch let err {
                    print("\(Date()) GET '\(url)' Could not decode JSON: \(err)")
                    
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
    
    func walks(offset: Int = 0, position: CLLocationCoordinate2D?, completion: @escaping ([Api.WalksItem]?, Error?) -> Void) {
        var url = URLComponents(string: "\(Api.url)/walks")!
        url.queryItems = [URLQueryItem(name: "offset", value: String(offset))]
        if let pos = position {
            url.queryItems!.append(contentsOf: [URLQueryItem(name: "lat", value: String(pos.latitude)), URLQueryItem(name: "lng", value: String(pos.longitude))])
        }
        get(url.url!, completion: completion)
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
    
    func walk(_ uuid: String, completion: @escaping (Api.Walk?, Error?) -> Void) {
        let url = URLComponents(string: "\(Api.url)/walk/\(uuid)")!
        get(url.url!, completion: completion)
    }
    
    // MARK: - FRIENDS
    
    struct GroupInfo: Codable {
        let uuid: String
        let display: String
        let handles: [String]
    }
    
    struct Friends: Codable {
        let friends: [GroupInfo]
        let randoms: [String]
    }
    
    func friends(completion: @escaping (Api.Friends?, Error?) -> Void) {
        let url = URLComponents(string: "\(Api.url)/me/friends")!
        get(url.url!, completion: completion)
    }
    
    // MARK: - MESSAGE
    
    struct MessageInput: Codable {
        let text: String
        let handles: [String]
    }
    
    struct MessageInput2: Codable {
        let text: String
        let groupUuid: String
    }
    
    struct MessageOutput: Codable {}
    
    func message(_ text: String, _ recipients: [String], completion: @escaping (Api.MessageOutput?, Error?) -> Void) {
        let url = URLComponents(string: "\(Api.url)/message")!
        post(url.url!, payload: MessageInput(text: text, handles: recipients), completion: completion)
    }
    
    func message(_ text: String, groupUuid: String, completion: @escaping (Api.MessageOutput?, Error?) -> Void) {
        let url = URLComponents(string: "\(Api.url)/message")!
        post(url.url!, payload: MessageInput2(text: text, groupUuid: groupUuid), completion: completion)
    }
    
    // MARK: - GROUPS
    
    struct GroupsOutput: Codable {
        let groups: [GroupsOutputGroup]
    }
    
    struct GroupsOutputGroup: Codable {
        let uuid: String
        let display: String
        let handles: [String]
    }
    
    func groups(completion: @escaping (Api.GroupsOutput?, Error?) -> Void) {
        let url = URLComponents(string: "\(Api.url)/me/groups")!
        get(url.url!, completion: completion)
    }
    
    // MARK: - CHAT
    
    struct GroupOutput: Codable {
        let uuid: String
        let display: String
        let handles: [String]
        let messages: [GroupOutputMessage]
    }
    
    struct GroupOutputMessage: Codable {
        let text: String
        let createdAt: String
        let authorHandle: String
    }
    
    func group(_ uuid: String, offset: Int, completion: @escaping (Api.GroupOutput?, Error?) -> Void) {
        var url = URLComponents(string: "\(Api.url)/group/\(uuid)")!
        url.queryItems = [URLQueryItem(name: "offset", value: String(offset))]
        
        get(url.url!, completion: completion)
    }
}
