import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";

@Injectable({
  providedIn: 'root'
})
export class UserService {
  URL = "http://localhost:8080";
  constructor(private http:HttpClient) { }

  /*esto va según lo que se tenga en el servidor*/
  getIndex(){
    return this.http.get(`${this.URL}`)
  }

  /**Funciones del proyecto */
  compilar(json: any){
    return this.http.post(`${this.URL}/compilar`, json)
  }
}
