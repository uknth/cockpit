<?php
function post_to_go($type, $action, $data, $headers) {
   $fields = '';

   foreach($data as $key => $value) { 
      $fields .= $key . '=' . $value . '&'; 
   }

   rtrim($fields, '&');

   $post = curl_init();

   curl_setopt($post, CURLOPT_URL, trim(GO_PATH.$type.'/'.$action,'/') );
   curl_setopt($post, CURLOPT_POST, count($data));
   curl_setopt($post, CURLOPT_POSTFIELDS, $fields);
   curl_setopt($post, CURLOPT_RETURNTRANSFER, 1);
   curl_setopt($post, CURLOPT_HTTPHEADER, $headers);

   $result = curl_exec($post);

   $header = curl_getinfo ( $post );

   echo '<pre>'; print_r($header);echo '</pre>';

   curl_close($post);

   return $result;
}

function redirectIfSession($parent){
   if($parent->session->userdata('user_id') > 0)
      redirect(base_url());
}

function redirectIfNoSession($parent){
   //echo intval($parent->session->userdata('user_id')); exit;
   if(!intval($parent->session->userdata('user_id')))
      redirect(base_url().'login');
}